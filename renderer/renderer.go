package renderer

import (
	"fmt"
	"os"

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

	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	r.Window.GLSwap()

	// Shader stuff

	vertSource, err := os.ReadFile("./shader.vert")
	if err != nil {
		return nil, fmt.Errorf("Failed to open vertex shader src: %v", err)
	}

	fragSource, err := os.ReadFile("./shader.frag")
	if err != nil {
		return nil, fmt.Errorf("Failed to open fragment shader src: %v", err)
	}

	vertShader := compileShader(vertSource, gl.VERTEX_SHADER)
	fragShader := compileShader(fragSource, gl.FRAGMENT_SHADER)

	return r, nil
}

// compileShader compile and return the shader
func compileShader(source []byte, shaderType uint32) uint32 {
	log.Panicf("unimplemented")
	return 0
}

// Quit quit and close the renderer
func (r *Renderer) Quit()  {
	sdl.Quit()
	r.Window.Destroy()
	sdl.GLDeleteContext(r.GlContext)
}
