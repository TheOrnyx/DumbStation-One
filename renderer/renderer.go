package renderer

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	WIN_WIDTH = 1024
	WIN_HEIGHT = 512
)

type Renderer struct {
	Window *sdl.Window
	GlContext sdl.GLContext
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

	return r, nil
}

// Quit quit and close the renderer
func (r *Renderer) Quit()  {
	sdl.Quit()
	r.Window.Destroy()
	sdl.GLDeleteContext(r.GlContext)
}
