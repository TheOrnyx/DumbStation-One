package gpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/utils"
)

// The GPUSTAT register entries
//
// NOTE - idk if this was smart to seperate it but idk at this point
type GpuStat struct {
	pageBaseX           uint8         // Texture pages base X coordinate (64 byte increment) - (Bits 0-4)
	pageBaseY           uint8         // Texture pages base Y coordinate (256 line increment) - (Bit 5)
	semiTransparency    uint8         // Semi transparency (not really sure what it does tbh) - (Bits 5-6)
	textureDepth        TextureDepth  // The texture page color depth - (Bits 7-8)
	dithering           bool          // Enable dithering from 24-bit to 15-bit - (Bit 9)
	allowDrawToDisplay  bool          // allow drawing to display area - (Bit 10)
	forceSetMaskBit     bool          // Set Mask-Bit when drawing pixels - (Bit 11)
	checkMaskBeforeDraw bool          // Draw pixels, false=always, true=not to masked areas - (Bit 12)
	interlaceField      bool          // NOTE avoid calling field directyl! use method for this!! - (Bit 13)
	textureDisable      bool          // when true disable textures (NOTE no PS2's have 2mb vram) - (Bit 15)
	horizontalRes       HorizontalRes // combination of the horizontal resolution bits - (Bits 16-18)
	verticalRes         VerticalRes   // Vertical resolution (TODO - says smth about bit22) - (Bit 19)
	videoMode           VideoMode     // VideoMode Either PAL or NTSC - (Bit 20)
	displayDepth        DisplayDepth  // Display area color depth - (Bit 21)
	verticalInterlace   bool          // vertical interlate - (Bit 22)
	displayDisabled     bool          // when true display is disabled - (Bit 23)
	intRequest          bool          // true when interrupt is requested (or active not sure TODO) - (Bit 24)
	dataRequest         uint8         // TODO implement this with a method and stuff - (Bit 25)
	readyToRecvWord     bool          // Ready to receive Cmd Word - (Bit 26)
	readyToSendVram     bool          // Ready to send VRAM to CPU - (Bit 27)
	readyToRecvDMA      bool          // Ready to receive DMA block - (Bit 28)
	dmaDirection        DMADirection  // DMA Direction - (Bits 29-30)
	// NOTE - ignoring bit 31 for now cuz it seems confusing
}

// NewGPUStat create and return an initialized gpustat instace
func NewGPUStat() GpuStat {
	return GpuStat{
		textureDepth:    T4Bit,
		interlaceField:  true, // TODO - is this the same one as the 'field' in the guide??
		horizontalRes:   HResFromFields(0, 0),
		verticalRes:     Y240Lines,
		videoMode:       Ntsc,
		displayDepth:    D15Bit,
		displayDisabled: true,
		dmaDirection:    DirOff,
	}
}

// Status return uint32 representation of the status register
func (g *GpuStat) Status() uint32 {
	var r uint32 = 0

	r |= uint32(g.pageBaseX) << 0
	r |= uint32(g.pageBaseY << 4)
	r |= uint32(g.semiTransparency) << 5
	r |= uint32(g.textureDepth) << 7
	r |= utils.BoolToUint32(g.dithering) << 9
	r |= utils.BoolToUint32(g.allowDrawToDisplay) << 10
	r |= utils.BoolToUint32(g.forceSetMaskBit) << 11
	r |= utils.BoolToUint32(g.checkMaskBeforeDraw) << 12
	r |= utils.BoolToUint32(g.interlaceField) << 13
	// ignore bit 14
	r |= utils.BoolToUint32(g.textureDisable) << 15
	r |= g.horizontalRes.intoStatus()

	// FIXME - Temporarily killed as since we don't emulate bit31
	// setting vres properly causes the BIOS to lock
	// r |= uint32(g.verticalRes) << 19
	
	r |= uint32(g.videoMode) << 20
	r |= uint32(g.displayDepth) << 21
	r |= utils.BoolToUint32(g.verticalInterlace) << 22
	r |= utils.BoolToUint32(g.displayDisabled) << 23
	r |= utils.BoolToUint32(g.intRequest) << 24
	// r |= uint32(g.dataRequest) << 25

	// NOTE - unfinished so atm we'll just pretend GPU is always ready
	// r |= utils.BoolToUint32(g.readyToRecvWord) << 26
	// r |= utils.BoolToUint32(g.readyToSendVram) << 27
	// r |= utils.BoolToUint32(g.readyToRecvDMA) << 28
	r |= 1 << 26
	r |= 1 << 27
	r |= 1 << 28

	r |= uint32(g.dmaDirection) << 29
	r |= 0 << 31

	// do bit 25 shit here
	dmaReq := g.dataReq(r)
	r |= dmaReq << 25

	return r
}

// dataReq return data request based on some other stuff
func (g *GpuStat) dataReq(stat uint32) uint32 {
	switch g.dmaDirection {
	case DirOff:
		return 0
	case DirFifo:
		return 1
	case DirCPUToGP0:
		return (stat >> 28) & 1 // TODO - replace this with calls to bit28 in gpustat later
	case DirGPUReadToCPU:
		return (stat >> 27) & 1 // TODO - replace with call to bit27 later
	}

	log.Panicf("Unknown DMA Direction %v", g.dmaDirection)
	return 0
}

// softReset called by GP1(0x00), performs a soft reset on the gpustat
// NOTE - this should write the like value 0x14802000 into GPUSTAT
func (g *GpuStat) softReset() {
	// TODO - once again the guide sets bit15 as like texture disable and idk why
	g.pageBaseX = 0
	g.pageBaseY = 0
	g.semiTransparency = 0
	g.textureDepth = T4Bit
	g.dithering = false
	g.allowDrawToDisplay = false
	g.forceSetMaskBit = false
	g.checkMaskBeforeDraw = false
	g.interlaceField = true
	g.textureDisable = false
	g.horizontalRes = HResFromFields(0, 0)
	g.verticalRes = Y240Lines
	g.videoMode = Ntsc
	g.displayDepth = D15Bit
	g.verticalInterlace = false
	g.displayDisabled = true
	g.intRequest = false
	g.dataRequest = 0
	g.readyToRecvWord = true
	g.readyToSendVram = false
	g.readyToRecvDMA = true
	g.dmaDirection = DirOff
}
