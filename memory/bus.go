package memory

import (
	"fmt"

	"github.com/TheOrnyx/psx-go/log"
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
	// check that memory address isn't unaligned
	if addr % 4 != 0 {
		return 0xF, fmt.Errorf("Unaligned load32 address: 0x%x\n", addr)
	}

	if offset, contains := BIOS_RANGE.Contains(addr); contains {
		return b.bios.load32(offset), nil
	}
	
	return 0xF, fmt.Errorf("Unknown load32 at address %v", addr)
	
}

// Store32 Store 32 bit value val in address addr
//
// TODO - maybe clean this up, it's kinda gross
func (b *Bus) Store32(addr, val uint32) error {
	// check that memory address isn't unaligned
	if addr % 4 != 0 {
		return fmt.Errorf("Unaligned load32 address: 0x%x, val:0x%x\n", addr, val)
	}

	if offset, contains := MEM_CONTROL.Contains(addr); contains {
		switch offset {
		case 0: // expansion 1 base address
			if val != 0x1f000000 {
				return fmt.Errorf("Bad expansion 1 base address: 0x%x", val)
			}

		case 4: // expansion 2 base address
			if val != 0x1f802000 {
				return fmt.Errorf("Bad expansion 2 base address: 0x%x", val)
			}

		default:
			log.Warnf("Unhandled write to MEM_CONTROL register")
		}

		return nil
	}

	if _, contains := RAM_SIZE.Contains(addr); contains {
		return nil // do nothing
	}

	
	
	return fmt.Errorf("Haven't implemented writing to address 0x%08x with val 0x%08x", addr, val)
}
