package renderer

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/TheOrnyx/psx-go/log"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIN_WIDTH = 1024
	WIN_HEIGHT = 512
)

type Renderer struct {
	Window *sdl.Window
	GlContext sdl.GLContext

	vertexShader uint32 // the Vertex shader object
	fragmentShader uint32 // The Fragment shader object
	program uint32 // OpenGL Program object
	vertexArrayObject uint32 // Vertex Array Object VAO
	positions Buffer[VRAMPos] // Buffer containing vertex positions
	colors Buffer[Color] // Buffer containing vertex colors
	numVertices uint32 // Current number of vertices in the buffers
	uniformOffset int32 // Index of the "offset" shader uniform
}

// NewRenderer create and initialize a new renderer object
func NewRenderer() (*Renderer, error) {
	r := new(Renderer)
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize SDL: %v", err)
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_DEBUG_FLAG)

	window, err := sdl.CreateWindow("PSX-GO", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_OPENGL)
	if err != nil {
		r.Quit()
		return nil, fmt.Errorf("Failed to create window: %v", err)
	}
	r.Window = window

	glContext, err := r.Window.GLCreateContext()
	if err != nil {
		r.Quit()
		return nil, fmt.Errorf("Failed to create gl context: %v", err)
	}
	r.GlContext = glContext

	if err := gl.Init(); err != nil {
		r.Quit()
		return nil, fmt.Errorf("Failed to initialize OpenGL: %v", err)
	}

	// enable debug output
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(DebugCallback, nil)

	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	r.Window.GLSwap()

	// Shader stuff

	vertSource, err := os.ReadFile("./renderer/shader.vert")
	if err != nil {
		return nil, fmt.Errorf("Failed to open vertex shader src: %v", err)
	}

	fragSource, err := os.ReadFile("./renderer/shader.frag")
	if err != nil {
		return nil, fmt.Errorf("Failed to open fragment shader src: %v", err)
	}

	vertShader := compileShader(vertSource, gl.VERTEX_SHADER)
	fragShader := compileShader(fragSource, gl.FRAGMENT_SHADER)

	// Link the program and then use it
	program := linkProgram([]uint32{vertShader, fragShader})
	gl.UseProgram(program)

	// generate vertex attribute object to hold vertex attributes
	var vao uint32 = 0
	gl.GenVertexArrays(1, &vao)
	// bind the VAO
	gl.BindVertexArray(vao)

	// Setup position attribute
	// Create buffer holding the positions
	positions := NewBuffer[VRAMPos]()

	// retrieve index for the attribute in shader and enable it
	index := findProgramAttrib(program, "vertex_position")
	gl.EnableVertexAttribArray(index)

	// link buffer and index: 2 non-normalized GLshort attributes
	gl.VertexAttribIPointer(index, 2, gl.SHORT, 0, nil)

	// Color stuff
	// Setup color attribute and bind it
	colors := NewBuffer[Color]()

	index = findProgramAttrib(program, "vertex_color")
	gl.EnableVertexAttribArray(index)

	// Link buffer and the index: 3 non-normalized GLByte
	// attributes. Should send data untouched to vertex shader
	gl.VertexAttribIPointer(index, 3, gl.UNSIGNED_BYTE, 0, nil)

	uniformOffset := gl.GetUniformLocation(program, gl.Str("offset"+"\x00")) // TODO - check
	gl.Uniform2i(uniformOffset, 0, 0)

	r.vertexShader = vertShader
	r.fragmentShader = fragShader
	r.program = program
	r.vertexArrayObject = vao
	r.positions = positions
	r.colors = colors
	r.numVertices = 0
	r.uniformOffset = uniformOffset
	
	return r, nil
}

// compileShader compile and return the shader
func compileShader(source []byte, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	// Attempt to compile shader
	cStr, free := gl.Strs(string(source)+"\x00")
	defer free()

	gl.ShaderSource(shader, 1, cStr, nil)
	gl.CompileShader(shader)

	// Extra bit of error checking in case we're not using debug
	// opengl context
	status := int32(gl.FALSE)
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status != gl.TRUE {
		log.Panic("Shader compilation failed!")
	}
	
	return shader
}

// linkProgram Link the shaders to the program
func linkProgram(shaders []uint32) uint32 {
	var program uint32

	program = gl.CreateProgram()

	for i := range shaders {
		gl.AttachShader(program, shaders[i])
	}

	gl.LinkProgram(program)

	var status int32 = gl.FALSE
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status != gl.TRUE {
		log.Panic("Opengl program linking failed!")
	}
	
	return program
}

// findProgramAttrib return index of attribute attribute in program
func findProgramAttrib(program uint32, attr string) uint32 {
	index := gl.GetAttribLocation(program, gl.Str(attr+"\x00"))

	if index < 0 {
		log.Panicf("Attribute %s, not found in program", attr)
	}

	return uint32(index)
}

// PushTriangle Add a triangle to the draw buffer
func (r *Renderer) PushTriangle(positions [3]VRAMPos, colors [3]Color)  {
	// make sure we have enough room left to queue the vertex
	if r.numVertices + 3 > VERTEX_BUFFER_LEN {
		log.Info("Vertex attrivute buffers full, forcing draw")
		r.Draw()
	}

	for i := range 3 {
		r.positions.Set(r.numVertices, positions[i])
		r.colors.Set(r.numVertices, colors[i])
		r.numVertices += 1
	}

}

// PushQuad Add a quad to the draw buffer
func (r *Renderer) PushQuad(positions [4]VRAMPos, colors [4]Color)  {
	// Make sure we have enough room left to queue the vertex.
	if r.numVertices + 6 > VERTEX_BUFFER_LEN {
		// Vertex attribute buffers are full, force early draw
		r.Draw()
	}

	// Push first triangle
	for i := range 3 {
		r.positions.Set(r.numVertices, positions[i])
		r.colors.Set(r.numVertices, colors[i])
		r.numVertices+=1
	}

	// Push second triangle
	for i := 1; i < 4; i++{
		r.positions.Set(r.numVertices, positions[i])
		r.colors.Set(r.numVertices, colors[i])
		r.numVertices+=1
	}
}

// Draw draw the buffered commands and reset the buffers
//
// TODO - improve later by using double buffering as this stalls the emulator
func (r *Renderer) Draw()  {
	// Make sure all the data from persisent mappings is flushed to
	// the buffer
	gl.MemoryBarrier(gl.CLIENT_MAPPED_BUFFER_BARRIER_BIT)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(r.numVertices))

	// Wait for GPU to complete
	sync := gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0)

	for  {
		r := gl.ClientWaitSync(sync, gl.SYNC_FLUSH_COMMANDS_BIT, 10000000)

		if r == gl.ALREADY_SIGNALED || r == gl.CONDITION_SATISFIED {
			// drawing is finished
			break
		}
	}

	// reset buffers
	r.numVertices = 0
}

// Display Draw the buffered commands and display them
func (r *Renderer) Display()  {
	r.Draw()
	r.Window.GLSwap()
}

// Quit quit and close the renderer
func (r *Renderer) Quit()  {
	gl.DeleteVertexArrays(1, &r.vertexArrayObject)
	gl.DeleteShader(r.vertexShader)
	gl.DeleteShader(r.fragmentShader)
	gl.DeleteProgram(r.program)
	
	sdl.Quit()
	r.Window.Destroy()
	sdl.GLDeleteContext(r.GlContext)
}

// DebugCallback Debug callback function for OpenGL
func DebugCallback(source, glType, id, severity uint32, length int32, msg string, userParam unsafe.Pointer)  {
	log.Infof("[OpenGL Debug] Source: 0x%x, Type: 0x%x, ID: %d, Severity: 0x%x, Message: %s\n",
		source, glType, id, severity, msg)
}

// SetDrawOffset Set value of the uniform draw offset
func (r *Renderer) SetDrawOffset(x, y int16)  {
	// Force draw for the primitives with the current offset
	r.Draw()

	gl.Uniform2i(r.uniformOffset, int32(x), int32(y))
}
