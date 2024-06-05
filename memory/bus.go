package memory

import (
	"fmt"
)

// the memory bus
type Bus struct {
	bios *Bios
}

// NewBus create and return a new bus object
func NewBus(bios *Bios) *Bus {
	return &Bus{bios: bios}
}

// Load32 load and return the value at addr on the bus
func (b *Bus) Load32(addr uint32) (uint32, error) {
	switch {
	case addr >= BIOS_LOWER && addr < BIOS_UPPER: // the bios memory range
		return b.bios.load32(addr - BIOS_LOWER), nil

	default:
		return 0xF, fmt.Errorf("Unknown load32 at address %v", addr)
	}
}
