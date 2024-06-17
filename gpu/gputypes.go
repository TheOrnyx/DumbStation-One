package gpu

import "github.com/TheOrnyx/psx-go/log"

// NOTE - i kinda just made this it's own file cuz I can't be bothered
// having half of my other files be filled up with a bunch of consts
// and shit (go please add enums )


type TextureDepth uint8

const (
	T4Bit  TextureDepth = 0
	T8Bit  TextureDepth = 1
	T15Bit TextureDepth = 2
)

// texDepthFromU32 get texturedepth value from uint32
func texDepthFromU32(val uint32) TextureDepth {
	switch val {
	case 0:
		return T4Bit
	case 1:
		return T8Bit
	case 2:
		return T15Bit
	default:
		log.Panicf("Failed to decode textureDepth from uint32, val:%v", val)
	}
	
	return 0 // shouldn't happen
}

type DisplayDepth uint8

// Display color depth constants
const (
	D15Bit DisplayDepth = 0
	D24Bit DisplayDepth = 1
)

type VerticalRes uint8

// Vertical Resolution constants
const (
	Y240Lines VerticalRes = 0
	Y480Lines VerticalRes = 1
)

type VideoMode uint8

// videoMode constants
const (
	Ntsc VideoMode = 0
	Pal  VideoMode = 1
)

// horizontal resolution type
type HorizontalRes uint8

// HResFromFields create and return a new HRes byte from the two fields
func HResFromFields(hr1, hr2 uint8) HorizontalRes {
	hr := (hr2 & 1) | ((hr1 & 3) << 1)

	return HorizontalRes(hr)
}

// intoStatus retrieve value of the bits for GpuStat for this
// horizontalRes instance
func (h HorizontalRes) intoStatus() uint32 {
	return uint32(h) << 16
}

type DMADirection uint8

// DMA direction constants
const (
	DirOff          DMADirection = 0
	DirFifo         DMADirection = 1 // NOTE - guide states this as FIFO but specs say "?"
	DirCPUToGP0     DMADirection = 2
	DirGPUReadToCPU DMADirection = 3 // NOTE - guide treats this as VRAMToCPU
)

