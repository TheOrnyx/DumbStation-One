package gpu

import "github.com/TheOrnyx/psx-go/log"

type Gpu struct {
	gpuStat GpuStat

	rectangleTextureXFlip bool // mirror textured rectangles along x axis
	rectangleTextureYFlip bool // Mirror textured rectangles along y axis
}

// NewGPU create and return a new gpu
func NewGPU() Gpu {
	return Gpu{gpuStat: NewGPUStat()}
}

// GP0 handle writes to the GP0 command register
func (g *Gpu) GP0(val uint32)  {
	opcode := (val >> 24) & 0xff // top byte contains opcode
	switch opcode {
	case 0xe1: // draw Mode
		g.gp0DrawMode(val)
	default:
		log.Panicf("Unhandled GP0 command: 0x%08x", val)
	}
}

// gp0DrawMode GPO(0xE1) command for setting draw mode settings
func (g *Gpu) gp0DrawMode(val uint32)  {
	stat := &g.gpuStat // can't be bothered typing g.gpuStat each time

	stat.pageBaseX = uint8(val & 0xf)
	stat.pageBaseY = uint8((val >> 4) & 1)
	stat.semiTransparency = uint8((val >> 5) & 3)

	stat.textureDepth = texDepthFromU32((val >> 7) & 3)

	stat.dithering = ((val >> 9) & 1) != 0
	stat.allowDrawToDisplay = ((val >> 10) & 1) != 0
	stat.useMaskForDrawing = ((val >> 11) & 1) != 0 // FIXME - the specs says this bit should set base y page 2

	g.rectangleTextureXFlip = ((val >> 12) & 1) != 0
	g.rectangleTextureYFlip = ((val >> 13) & 1) != 0
}
