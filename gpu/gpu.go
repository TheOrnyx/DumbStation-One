package gpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/renderer"
)

// The GPU Struct
//
// TODO - rename these fields, they're really bad names too long
// TODO - maybe like seperate GP0 and GP1 stuff into their own type idk
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

	gp0CmdBuffer      CommandBuffer // Buffer containing current GP0 command
	gp0WordsRemaining uint32        // The remaining words for the current GP0 command
	gp0Cmd            GP0Cmd        // the GPU command for holding the length, function etc
	gp0Mode           GP0Mode       // The current mode of the GP0 register

	renderer *renderer.Renderer // The OpenGL Renderer
}

// NewGPU create and return a new gpu
func NewGPU() Gpu {
	var err error
	g := Gpu{gpuStat: NewGPUStat(), gp0Mode: GP0ModeCommand}
	g.renderer, err = renderer.NewRenderer()
	if err != nil {
		log.Panicf("Failed to create renderer: %v", err)
	}

	return g
}

// Status return the status register
func (g *Gpu) Status() uint32 {
	return g.gpuStat.Status()
}

// GP0 handle writes to the GP0 command register
func (g *Gpu) GP0(val uint32) {
	if g.gp0WordsRemaining == 0 {
		opcode := (val >> 24) & 0xff // top byte contains opcode

		cmd, valid := gp0Commands[opcode]
		if !valid {
			log.Panicf("Unhandled GP0 command: 0x%08x, Opcode:0x%02x", val, opcode)
		}

		g.gp0WordsRemaining = cmd.length
		g.gp0Cmd = cmd
		g.gp0CmdBuffer.clear()
	}

	g.gp0WordsRemaining -= 1

	switch g.gp0Mode {
	case GP0ModeCommand:
		g.gp0CmdBuffer.pushWord(val)

		if g.gp0WordsRemaining == 0 {
			// we have alll parameters, we can run command
			g.gp0Cmd.runFunc(g, val)
		}

	case GP0ModeImgLoad:
		// TODO - should copy pixel data to vram
		if g.gp0WordsRemaining == 0 {
			g.gp0Mode = GP0ModeCommand
		}

	default:
		log.Panicf("Unknown GP0 command")
	}
}

// GP1 Handle writes to the GP1 command register
func (g *Gpu) GP1(val uint32) {
	opcode := (val >> 24) & 0xff

	switch opcode {
	case 0x00: // GP1 reset
		g.gp1Reset(val)
	case 0x01:
		g.gp1ResetCmdBuffer()
	case 0x02:
		g.gp1AcknowledgeInterrupt()
	case 0x03:
		g.gp1DisplayEnable(val)
	case 0x04: // GP1 DMA direction
		g.gp1SetDMADirection(val)
	case 0x05:
		g.gp1SetDisplayVRAMStart(val)
	case 0x06:
		g.gp1SetHorizDisplayRange(val)
	case 0x07:
		g.gp1SetVertDisplayRange(val)
	case 0x08: // GP1 display mode
		g.gp1DisplayMode(val)
	default:
		log.Panicf("Unhandled GP1 command: 0x%08x, Opcode:0x%02x", val, opcode)
	}
}

// gp1Reset GP1(0x00): Soft reset of the GPU
func (g *Gpu) gp1Reset(val uint32) {
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
	g.displayHorizEnd = 0x200 + 256*10

	// Reset vertical display range GP1(07h)
	g.displayLineStart = 0x010
	g.displayLineEnd = 0x010 + 240

	g.gpuStat.softReset() // TODO - check
	g.gp1ResetCmdBuffer()
}

// Read retrieve the value of the read register
func (g *Gpu) Read() uint32 {
	log.Info("Not implemented GPUREAD yet")
	return 0
}

// gp1DisplayMode GP1(08h) - Display mode
func (g *Gpu) gp1DisplayMode(val uint32) {
	stat := &g.gpuStat
	hr1 := uint8(val & 3)
	hr2 := uint8((val >> 6) & 1)
	stat.horizontalRes = HResFromFields(hr1, hr2)

	stat.verticalRes = Y240Lines
	if val&0x4 != 0 {
		stat.verticalRes = Y480Lines
	}

	stat.videoMode = Ntsc
	if val&0x8 != 0 {
		stat.videoMode = Pal
	}

	stat.displayDepth = D24Bit
	if val&0x10 != 0 {
		stat.displayDepth = D15Bit
	}

	stat.verticalInterlace = val&0x20 != 0

	if val&0x80 != 0 { // flipped
		log.Panicf("Horizontal flip not implemented yet!")
	}
}

// gp1SetDMADirection GP1(04h) - Set DMA Direction
func (g *Gpu) gp1SetDMADirection(val uint32) {
	switch val & 3 {
	case 0:
		g.gpuStat.dmaDirection = DirOff
	case 1:
		g.gpuStat.dmaDirection = DirFifo
	case 2:
		g.gpuStat.dmaDirection = DirCPUToGP0
	case 3:
		g.gpuStat.dmaDirection = DirGPUReadToCPU
	default:
		log.Panicf("Unknown DMA direction %v", val&3)
	}
}

// gp1SetDisplayVRAMStart GP1(05h) - Set Display VRAM start address
func (g *Gpu) gp1SetDisplayVRAMStart(val uint32) {
	g.displayVramXStart = uint16(val & 0x3fe)
	g.displayVramYStart = uint16((val >> 10) & 0x1ff)
}

// gp1SetHorizDisplayRange GP1(06h) - Set Horizontal display range on screen
func (g *Gpu) gp1SetHorizDisplayRange(val uint32) {
	g.displayHorizStart = uint16(val & 0xfff)
	g.displayHorizEnd = uint16((val >> 12) & 0xfff)
}

// gp1SetVertDisplayRange GP1(07h) - Set vertical display range on screen
func (g *Gpu) gp1SetVertDisplayRange(val uint32) {
	g.displayLineStart = uint16(val & 0x3ff)
	g.displayLineEnd = uint16((val >> 10) & 0x3ff)
}

// gp1DisplayEnable GP1(03h) - Display Enable
func (g *Gpu) gp1DisplayEnable(val uint32) {
	g.gpuStat.displayDisabled = val&1 != 0
}

// gp1AcknowledgeInterrupt GP1(02h) - Acknowledge GPU Interrupt (IRQ1)
func (g *Gpu) gp1AcknowledgeInterrupt() {
	g.gpuStat.intRequest = false
}

// gp1ResetCmdBuffer GP1(01h) - Reset Command Buffer
func (g *Gpu) gp1ResetCmdBuffer() {
	g.gp0CmdBuffer.clear()
	g.gp0WordsRemaining = 0
	g.gp0Mode = GP0ModeCommand
	// TODO - should clear command FIFO when implemented
}
