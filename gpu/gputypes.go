package gpu

import (
	"unsafe"

	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

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

// Type for states for the Gp0 command register
type GP0Mode uint8

const (
	GP0ModeCommand GP0Mode = 0 // Default mode: Handling commands
	GP0ModeImgLoad GP0Mode = 1 // Loading image into VRAM
)

// Position in VRAM
type VRAMPos struct {
	x int16
	y int16
}

// PosFromGP0 parse vram position from a GP0 param
func PosFromGP0(val uint32) VRAMPos {
	x := int16(val)
	y := int16(val >> 16)

	return VRAMPos{x: x, y: y}
}

// RGB color
type Color struct {
	r uint8
	g uint8
	b uint8
}

// ColorFromGP0 Parse color from a GP0 param
func ColorFromGP0(val uint32) Color {
	r := uint8(val)
	g := uint8(val >> 8)
	b := uint8(val >> 16)

	return Color{r, g, b}
}

const VERTEX_BUFFER_LEN uint32 = 64*1024

type Buffer[T any] struct {
	object uint32
	mapBuffer *T
}

// NewBuffer create and return a new Buffer object
func NewBuffer[T any]() Buffer[T] {
	var object uint32
	var memory *T

	// generate buffer object
	gl.GenBuffers(1, &object)

	// bind it
	gl.BindBuffer(gl.ARRAY_BUFFER, object)

	// compute size of buffer
	elemSize := int(unsafe.Sizeof(memory))
	bufferSize := elemSize * int(VERTEX_BUFFER_LEN)

	// Write only persisent mapping (NOTE - guide says not coherent)
	access := gl.MAP_WRITE_BIT | gl.MAP_PERSISTENT_BIT

	// Allocate buffer memory
	gl.BufferStorage(gl.ARRAY_BUFFER, bufferSize, nil, uint32(access))

	// Remap entire buffer
	memory = (*T)(gl.MapBufferRange(gl.ARRAY_BUFFER, 0, bufferSize, uint32(access)))

	// Reset buffer to 0 to avoid hard to reproduce bugs in case we do
	// something wrong with uniitiliazed memory
	s := unsafe.Slice(memory, VERTEX_BUFFER_LEN)
	for i := range s {
		var r T
		s[i] = r
	}
	
	return Buffer[T]{
		object: object,
		mapBuffer: memory,
	}
}

// Set Set entry at index to val in the buffer
func (b *Buffer[T]) Set(index uint32, val T)  {
	if index >= VERTEX_BUFFER_LEN {
		log.Panicf("Buffer overflow!")
	}

	elemSize := unsafe.Sizeof(*b.mapBuffer)
	p := (*T)(unsafe.Pointer(uintptr(unsafe.Pointer(b.mapBuffer)) + uintptr(index) * elemSize))
	*p = val
}

// Clear clear the buffer objects
func (b *Buffer[T]) Clear()  {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.object)
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
	gl.DeleteBuffers(1, &b.object)
}
