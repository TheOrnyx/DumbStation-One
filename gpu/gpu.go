package gpu

import "github.com/TheOrnyx/psx-go/log"

// The GPU Struct
//
// TODO - rename these fields, they're really bad names too long
type Gpu struct {
	gpuStat GpuStat

	rectangleTextureXFlip bool   // mirror textured rectangles along x axis
	rectangleTextureYFlip bool   // Mirror textured rectangles along y axis
	texWindowXMask        uint8  // texture window x mask (8 pixel steps)
	texWindowYMask        uint8  // texture window y mask (8 pixel steps)
	texWindowXOffset      uint8  // texture window x offset (8 pixel steps)
	texWindowYOffset      uint8  // texture window y offset (8 pixel steps)
	drawAreaLeft          uint16 // Left-most column of drawing area
	drawAreaTop           uint16 // Top-most line of drawing area
	drawAreaRight         uint16 // Right-most column of drawing area
	drawAreaBottom        uint16 // Bottom-most line of drawing area
	drawXOffset           int16  // Horizontal drawing offset applied to all the vertex
	drawYOffset           int16  // Vertical drawing offset applied to all the vertex :D
	displayVramXStart     uint16 // First column of the display area in VRAM
	displayVramYStart     uint16 // First line of the display area in VRAM
	displayHorizStart     uint16 // Display output horizontal start relative to HSYNC 
	displayHorizEnd       uint16 // Display output horizontal end relative to HSYNC
	displayLineStart      uint16 // Display output first line relative to VSYNC
	displayLineEnd        uint16 // Display output last line relative to VSYNC
}

// NewGPU create and return a new gpu
func NewGPU() Gpu {
	return Gpu{gpuStat: NewGPUStat()}
}

// GP0 handle writes to the GP0 command register
func (g *Gpu) GP0(val uint32) {
	opcode := (val >> 24) & 0xff // top byte contains opcode
	switch opcode {
	case 0x00: // NOP
		break
	case 0xe1: // draw Mode
		g.gp0DrawMode(val)
	default:
		log.Panicf("Unhandled GP0 command: 0x%08x", val)
	}
}

// gp0DrawMode GPO(0xE1) command for setting draw mode settings
func (g *Gpu) gp0DrawMode(val uint32) {
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


// GP1 Handle writes to the GP1 command register
func (g *Gpu) GP1(val uint32)  {
	opcode := (val >> 24) & 0xff

	switch opcode {
	case 0x00: // GP1 reset
		g.gp1Reset(val)
	default:
		log.Panicf("Unhandled GP1 command: 0x%08x", val)
	}
}

// gp1Reset GP1(0x00): Soft reset of the GPU
func (g *Gpu) gp1Reset(val uint32)  {
	// TODO - this should clear the command FIFO later
	// TODO - this should invalidate GPU cache if it ever gets implemented
	// TODO - check this, kinda winging it based on the stuff

	
	g.rectangleTextureXFlip = false
	g.rectangleTextureYFlip = false

	// Reset vRAM display stuff?? GP1(05h)
	// TODO - check this
	g.displayVramXStart = 0
	g.displayVramYStart = 0
	
	// Reset horizontal display range GP1(06h)
	// TODO - check this is correct
	g.displayHorizStart = 0x200
	g.displayHorizEnd = 0x200 + 256 * 10

	// Reset vertical display range GP1(07h)
	g.displayLineStart = 0x010
	g.displayLineEnd = 0x010+240
	
	g.gpuStat.softReset() // TODO - check
}
