package memory

import (
	"fmt"
	"os"

	"github.com/TheOrnyx/psx-go/utils"
)

const BIOS_SIZE = 0x80000 // the size of the bios (512KB)
const (
	BIOS_LOWER = 0xbfc00000 // the start of the bios memory
	BIOS_UPPER = BIOS_LOWER + BIOS_SIZE // the end of the bios memory
)

type Bios struct {
	data []uint8 // the bios data
}

// NewBios create and return a new bios object from a path
// TODO - check if returning a pointer is slower or not
func NewBios(path string) (*Bios, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read BIOS file: %v", err)
	}

	if len(data) == BIOS_SIZE {
		return &Bios{data: data}, nil
	}

	return nil, fmt.Errorf("data was not expected bios size: Expected:%d, Received:%d", BIOS_SIZE, len(data))
}

// load32 get and return the 32bit little endian word at 'offset'
// where offset is an offset from the start of the bios data, not the absolute address
func (b *Bios) load32(offset uint32) uint32 {
	b0 := b.data[offset + 0]
	b1 := b.data[offset + 1]
	b2 := b.data[offset + 2]
	b3 := b.data[offset + 3]

	return utils.BytesToUint32(b0,b1,b2,b3)
}
