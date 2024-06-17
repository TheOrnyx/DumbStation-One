package gpu

// This is for like commands and shit for GP0 cuz cbf putting it in the same file

// GP0 handle writes to the GP0 command register
func (g *Gpu) GP0(val uint32) {
	opcode := (val >> 24) & 0xff // top byte contains opcode
	switch opcode {
	case 0x00: // NOP
		break
	case 0xe1: // draw Mode
		g.gp0DrawMode(val)
	case 0xe2:
		g.gp0SetTextureWindow(val)
	case 0xe3: // drawareatopleft
		g.gp0SetDrawAreaTopLeft(val)
	case 0xe4:
		g.gp0SetDrawAreaBtmRight(val)
	case 0xe5:
		g.gp0SetDrawOffset(val)
	case 0xe6:
		g.gp0SetMaskBitSetting(val)
	default:
		log.Panicf("Unhandled GP0 command: 0x%08x, Opcode:0x%02x", val, opcode)
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
	stat.forceSetMaskBit = ((val >> 11) & 1) != 0 // FIXME - the specs says this bit should set base y page 2

	g.rectangleTextureXFlip = ((val >> 12) & 1) != 0
	g.rectangleTextureYFlip = ((val >> 13) & 1) != 0
}

// gp0SetDrawAreaTopLeft GP0(E3h) - Set top left drawing area
func (g *Gpu) gp0SetDrawAreaTopLeft(val uint32) {
	g.drawAreaTop = uint16((val >> 10) & 0x3ff)
	g.drawAreaLeft = uint16(val & 0x3ff)
}

// gp0SetDrawAreaBtmRight GP0(E4h) - Set bottom right drawing area
func (g *Gpu) gp0SetDrawAreaBtmRight(val uint32) {
	g.drawAreaBottom = uint16((val >> 10) & 0x3ff)
	g.drawAreaRight = uint16(val & 0x3ff)
}

// gp0SetDrawOffset GP0(E5h) - Set Drawing Offset
func (g *Gpu) gp0SetDrawOffset(val uint32) {
	x := uint16(val & 0x7ff)
	y := uint16((val >> 11) & 0x7ff)

	// values are 11bit two's complement signed values so shift to force sign extension
	g.drawXOffset = int16(x<<5) >> 5
	g.drawYOffset = int16(y<<5) >> 5
}

// gp0SetTextureWindow GP0(E2h) - Set Texture Window
func (g *Gpu) gp0SetTextureWindow(val uint32) {
	g.texWindowXMask = uint8(val & 0x1f)
	g.texWindowYMask = uint8((val >> 5) & 0x1f)
	g.texWindowXOffset = uint8((val >> 10) & 0x1f)
	g.texWindowYOffset = uint8((val >> 15) & 0x1f)
}

// gp0SetMaskBitSetting GP0(E6h) - Set mask bit settings
func (g *Gpu) gp0SetMaskBitSetting(val uint32) {
	g.gpuStat.forceSetMaskBit = val&1 != 0
	g.gpuStat.checkMaskBeforeDraw = val&2 != 0
}



//////////////////
// Command list //
//////////////////

// GP0 command struct
type GP0Cmd struct {
	opcode uint32 // the opcode
	name string // the name
	length uint8 // The number of commands for this command
	runFunc func(g *Gpu, val uint32) // the run function
}

var gp0Commands map[uint32]GP0Cmd = map[uint32]GP0Cmd{
	0x00: GP0Cmd{0x00, "NOP", func(g *Gpu, val uint32) {}}
}
