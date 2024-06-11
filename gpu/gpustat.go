package gpu

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
	flipHorizontally   bool          // Flip screen horizontally (somth about v1 only) - (Bit 14)
	pageBaseY2         uint8         // Texture page y Base 2 (only for 2MB VRAM) - (Bit 15)
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
