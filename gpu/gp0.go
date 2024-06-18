package gpu

import "github.com/TheOrnyx/psx-go/log"

// This is for like commands and shit for GP0 cuz cbf putting it in the same file



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

// gp0Nop GP0(00h) Nop
func (g *Gpu) gp0Nop() {
	// nothing :3
}

// gp0ClearCache GP0(01h) - Clear Texture cache
func (g *Gpu) gp0ClearCache()  {
	log.Info("Texture cache not implemented yet")
}

// gp0MonoQuadPolyOpaque GP0(28h) - Monochrome four-point polygon, opaque
func (g *Gpu) gp0MonoQuadPolyOpaque(val uint32) {
	log.Infof("(not implemented yet) Draw quad :3")
}

// gp0ImageLoad GP0(A0h) - Image load
func (g *Gpu) gp0ImageLoad()  {
	// param 2 contains image res
	res := g.gp0CmdBuffer.at(2)

	width := res & 0xffff
	height := res >> 16
	imgSize := width*height

	// If we have off number of pixels we must round up since we
	// transfer 32bits at a time there'll be 16-bits padding in the
	// last word
	imgSize = (imgSize + 1) & (^ uint32(1))

	// Store number of words expected for image
	g.gp0WordsRemaining = imgSize / 2
	g.gp0Mode = GP0ModeImgLoad
}

// gp0ImageStore GP0(C0h) - Image Store
func (g *Gpu) gp0ImageStore()  {
	res := g.gp0CmdBuffer.at(2)

	width := res & 0xffff
	height := res >> 16

	log.Infof("(Not implemented yet) Unhandled image store: width:%v, height:%v", width, height)
}

// gp0QuadShadedOpaque GP0(38h) - Shaded opaque Quadrilateral
func (g *Gpu) gp0QuadShadedOpaque()  {
	log.Info("(Not implemented yet) Draw quad shaded opaque")
}

// gp0TriShadedOpaque GP0(30h) - Shaded three-point polygon, opaque
func (g *Gpu) gp0TriShadedOpaque()  {
	log.Info("(Not implemented yet) Draw triangle shaded opaque")
}

// gp0QuadBlendedOpaque GP0(2Ch) - Textured four-point polygon, opaque, texture-blending
func (g *Gpu) gp0QuadBlendedOpaque()  {
	log.Info("(Not implemented yet) Draw texture-blended opaque quad")
}

//////////////////
// Command list //
//////////////////

// GP0 command struct
type GP0Cmd struct {
	opcode  uint32                   // the opcode
	length  uint32                    // The number of commands for this command
	name    string                   // the name
	runFunc func(g *Gpu, val uint32) // the run function
}

var gp0Commands map[uint32]GP0Cmd = map[uint32]GP0Cmd{
	0x00: {0x00, 1, "NOP", func(g *Gpu, val uint32) { g.gp0Nop() }},
	0x01: {0x01, 1, "Clear Cache", func(g *Gpu, val uint32) { g.gp0ClearCache() }},
	0x28: {0x28, 5, "Monochrome four-point polygon, opaque", func(g *Gpu, val uint32) { g.gp0MonoQuadPolyOpaque(val) }},
	0x2c: {0x2c, 9, "Textured four-point polygon, opaque, texture-blending", func(g *Gpu, val uint32) { g.gp0QuadBlendedOpaque() }},
	0x30: {0x30, 6, "Shaded three-point polygon, opaque", func(g *Gpu, val uint32) { g.gp0TriShadedOpaque() }},
	0x38: {0x38, 8, "Shaded four-point polygon, opaque", func(g *Gpu, val uint32) { g.gp0QuadShadedOpaque() }},
	0xa0: {0xa0, 3, "GP0 Image Load", func(g *Gpu, val uint32) { g.gp0ImageLoad() }},
	0xc0: {0xc0, 3, "Copy Rectangle (VRAM to CPU)/Image store", func(g *Gpu, val uint32) { g.gp0ImageStore() }},
	0xe1: {0xe1, 1, "Draw Mode setting", func(g *Gpu, val uint32) { g.gp0DrawMode(val) }},
	0xe2: {0xe2, 1, "Set Texture Window", func(g *Gpu, val uint32) { g.gp0SetTextureWindow(val) }},
	0xe3: {0xe3, 1, "Set Draw Area Top Left", func(g *Gpu, val uint32) { g.gp0SetDrawAreaTopLeft(val) }},
	0xe4: {0xe4, 1, "Set Draw Area Bottom Right", func(g *Gpu, val uint32) { g.gp0SetDrawAreaBtmRight(val) }},
	0xe5: {0xe5, 1, "Set Draw Offset", func(g *Gpu, val uint32) { g.gp0SetDrawOffset(val) }},
	0xe6: {0xe6, 1, "Set Mask Bit Setting", func(g *Gpu, val uint32) { g.gp0SetMaskBitSetting(val) }},
}
