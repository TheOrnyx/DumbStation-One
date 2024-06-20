package renderer

import (
	"unsafe"

	"github.com/TheOrnyx/psx-go/log"
	"github.com/go-gl/gl/v3.3-core/gl"
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
