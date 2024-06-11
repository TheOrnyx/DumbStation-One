package gpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/utils"
)

// The GPUSTAT register entries
//
// NOTE - idk if this was smart to seperate it but idk at this point
type GpuStat struct {
	pageBaseX          uint8         // Texture pages base X coordinate (64 byte increment) - (Bits 0-4)
	pageBaseY          uint8         // Texture pages base Y coordinate (256 line increment) - (Bit 5)
	semiTransparency   uint8         // Semi transparency (not really sure what it does tbh) - (Bits 5-6)
	textureDepth       TextureDepth  // The texture page color depth - (Bits 7-8)
	dithering          bool          // Enable dithering from 24-bit to 15-bit - (Bit 9)
	allowDrawToDisplay bool          // allow drawing to display area - (Bit 10)
	useMaskForDrawing  bool          // Set Mask-Bit when drawing pixels - (Bit 11)
	avoidDrawOnMask    bool          // Draw pixels, false=always, true=not to masked areas - (Bit 12)
	interlaceField     bool          // NOTE avoid calling field directyl! use method for this!! - (Bit 13)
	// flipHorizontally   bool          // Flip screen horizontally (somth about v1 only) - (Bit 14)
	pageBaseY2         uint8         // (NOTE unsure) Texture page y Base 2 (only for 2MB VRAM) - (Bit 15)
	horizontalRes      HorizontalRes // combination of the horizontal resolution bits - (Bits 16-18)
	verticalRes        VerticalRes   // Vertical resolution (TODO - says smth about bit22) - (Bit 19)
	videoMode          VideoMode     // VideoMode Either PAL or NTSC - (Bit 20)
	displayDepth       DisplayDepth  // Display area color depth - (Bit 21)
	verticalInterlace  bool          // vertical interlate - (Bit 22)
	displayDisabled    bool          // when true display is disabled - (Bit 23)
	intRequest         bool          // true when interrupt is requested (or active not sure TODO) - (Bit 24)
	dataRequest        uint8         // TODO implement this with a method and stuff - (Bit 25)
	readyToRecvWord    bool          // Ready to receive Cmd Word - (Bit 26)
	readyToSendVram    bool          // Ready to send VRAM to CPU - (Bit 27)
	readyToRecvDMA     bool          // Ready to receive DMA block - (Bit 28)
	dmaDirection       DMADirection  // DMA Direction - (Bits 29-30)
	// NOTE - ignoring bit 31 for now cuz it seems confusing
}

// NewGPUStat create and return an initialized gpustat instace
func NewGPUStat() GpuStat {
	return GpuStat{
		textureDepth: T4Bit,
		interlaceField: true, // TODO - is this the same one as the 'field' in the guide??
		horizontalRes: HResFromFields(0, 0),
		verticalRes: Y240Lines,
		videoMode: Ntsc,
		displayDepth: D15Bit,
		displayDisabled: true,
		dmaDirection: DirOff,
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
	r |= utils.BoolToUint32(g.useMaskForDrawing) << 11
	r |= utils.BoolToUint32(g.avoidDrawOnMask) << 12
	r |= utils.BoolToUint32(g.interlaceField) << 13
	// ignore bit 14
	r |= uint32(g.pageBaseY2) << 15
	r |= g.horizontalRes.intoStatus()
	r |= uint32(g.verticalRes) << 19
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
