package memory

import (
	"fmt"

	"github.com/TheOrnyx/psx-go/log"
)

// the memory bus
type Bus struct {
	bios *Bios
	ram  Ram
	dma  Dma // the DMA registers
}

// NewBus create and return a new bus object
func NewBus(bios *Bios) *Bus {
	return &Bus{bios: bios, ram: NewRam(), dma: NewDMA()}
}

// ReadDMAReg read the dma register
func (b *Bus) ReadDMAReg(offset uint32) uint32 {
	major := (offset & 0x70) >> 4
	minor := offset & 0xf

	switch major {
	// per channel registers
	case 0,1,2,3,4,5,6: // TODO - this is gross, kill later cbf rn
		channel := b.dma.GetChannelRef(PortFromIndex(major))
		switch minor {
		case 8:
			return channel.Control()
		default:
			log.Panicf("unhandled DMA read at: 0x%08x", offset)
		}

	case 7: // Common DMA registers
		switch minor {
		case 0:
			return b.dma.Control()
		case 4:
			return b.dma.Interrupt()
		default:
			log.Panicf("unhandled DMA read at: 0x%08x", offset)
		}
	}

	log.Panicf("unhandled DMA read at: 0x%08x", offset)
	return 0x00 // shouldn't be reached
}

// SetDMAReg DMA register write
func (b *Bus) SetDMAReg(offset, val uint32) {
	major := (offset & 0x70) >> 4
	minor := offset & 0xf

	switch major {
	// per channel registers
	case 0,1,2,3,4,5,6: // TODO - this is gross, kill later cbf rn
		channel := b.dma.GetChannelRef(PortFromIndex(major))
		switch minor {
		case 8:
			channel.SetControl(val)
		default:
			log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x", val, offset)
		}

	case 7: // Common DMA registers
		switch minor {
		case 0:
			b.dma.SetControl(val)
		case 4:
			b.dma.SetInterrupt(val)
		default:
			log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x", val, offset)
		}

	default:
		log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x", val, offset)
	}
}

// Load32 load and return the value at addr on the bus
func (b *Bus) Load32(addr uint32) (uint32, error) {
	// check that memory address isn't unaligned
	// if addr % 4 != 0 {
	// 	return 0xF, fmt.Errorf("Unaligned load32 address: 0x%x\n", addr)
	// }

	absAddr := MaskRegion(addr)

	if offset, contains := BIOS_RANGE.Contains(absAddr); contains {
		return b.bios.load32(offset), nil
	}

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		return b.ram.load32(offset), nil
	}

	if _, contains := IRQ_CONTROL.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) IRQ control 32bit read at: 0x%08x", absAddr)
		return 0, nil
	}

	if offset, contains := DMA_RANGE.Contains(absAddr); contains {
		return b.ReadDMAReg(offset), nil
	}

	if offset, contains := GPU_RANGE.Contains(absAddr); contains {
		log.Infof("(Not fully implemented yet) GPU 32bit read at: 0x%08x", absAddr)
		switch offset {
		case 4: // gpustat set bit 28 so gpu is ready to receive DMA blocks
			return 0x10000000, nil
		default:
			return 0, nil
		}
	}

	return 0xF, fmt.Errorf("Unknown load32 at address 0x%08x", addr)
}

// Load16 load 16-bit halfword at addr
func (b *Bus) Load16(addr uint32) (uint16, error) {
	absAddr := MaskRegion(addr)

	if _, contains := SPU_RANGE.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) SPU register 16bit read at: 0x%08x", absAddr)
		return 0, nil
	}

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		return b.ram.load16(offset), nil
	}

	if _, contains := IRQ_CONTROL.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) IRQ control 16bit read at: 0x%08x", absAddr)
		return 0, nil
	}

	return 0, fmt.Errorf("Unkown Load16 at address 0x%08x", absAddr)
}

// Load8 load byte at given address addr
func (b *Bus) Load8(addr uint32) (uint8, error) {
	absAddr := MaskRegion(addr)

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		return b.ram.load8(offset), nil
	}

	if offset, contains := BIOS_RANGE.Contains(absAddr); contains {
		return b.bios.load8(offset), nil
	}

	if _, contains := EXPANSION_1.Contains(absAddr); contains {
		// not implemented
		// TODO - i am so confused here, figure it out lmao
		log.Infof("(Not implemented yet) Expansion 1 8bit read at: absAddr:0x%08x addr:0x%08x", absAddr, addr)
		return 0xff, nil
	}

	return 0xF, fmt.Errorf("Unkown Load8 at address 0x%08x", absAddr)
}

// Store32 Store 32 bit value val in address addr
//
// TODO - maybe clean this up, it's kinda gross
func (b *Bus) Store32(addr, val uint32) error {
	// check that memory address isn't unaligned
	// if addr % 4 != 0 {
	// 	return fmt.Errorf("Unaligned load32 address: 0x%x, val:0x%x\n", addr, val)
	// }

	absAddr := MaskRegion(addr)

	if offset, contains := SYS_CONTROL.Contains(absAddr); contains {
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
			log.Info("Unhandled write to MEM_CONTROL register")
		}

		return nil
	}

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		b.ram.store32(offset, val)
		return nil
	}

	if _, contains := RAM_SIZE.Contains(absAddr); contains {
		return nil // do nothing
	}

	if _, contains := CACHE_CONTROL.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) Cache 32bit write 0x%08x to 0x%08x", val, absAddr)
		return nil
	}

	if _, contains := IRQ_CONTROL.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) IRQ Control 32bit write 0x%08x to 0x%08x", val, absAddr)
		return nil
	}

	if offset, contains := DMA_RANGE.Contains(absAddr); contains {
		b.SetDMAReg(offset, val)
		return nil
	}

	if _, contains := GPU_RANGE.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) GPU 32bit write 0x%08x to 0x%08x", val, absAddr)
		return nil
	}

	if _, contains := TIMERS_RANGE.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) TIMERS 32bit write 0x%08x to 0x%08x", val, absAddr)
		return nil
	}

	return fmt.Errorf("Haven't implemented store32 to address 0x%08x with val 0x%08x", addr, val)
}

// Store16 store 16 bit value into memory
func (b *Bus) Store16(addr uint32, val uint16) error {
	// if addr % 2 != 0 {
	// 	return fmt.Errorf("Unaligned Store16 address: 0x%x, val:0x%x\n", addr, val)
	// }

	absAddr := MaskRegion(addr)

	if _, contains := SPU_RANGE.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) SPU register 16bit write 0x%04x to 0x%08x", val, absAddr)
		return nil
	}

	if _, contains := TIMERS_RANGE.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) TIMERS register 16bit write 0x%04x to 0x%08x", val, absAddr)
		return nil
	}

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		b.ram.store16(offset, val)
		return nil
	}

	if _, contains := IRQ_CONTROL.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) IRQ control 16-bit write 0x%04x to 0x%08x", val, absAddr)
		return nil
	}

	return fmt.Errorf("Haven't implemented store16 into address 0x%08x with val 0x%04x", addr, val)
}

// Store8 store 8 bit value into memory
func (b *Bus) Store8(addr uint32, val uint8) error {
	absAddr := MaskRegion(addr)

	if offset, contains := RAM_RANGE.Contains(absAddr); contains {
		b.ram.store8(offset, val)
		return nil
	}

	if offset, contains := EXPANSION_2.Contains(absAddr); contains {
		log.Infof("(Not implemented yet) EXPANSION_2 8bit write 0x%02x to offset:0x%08x, absAddr:0x%02x", val, offset, absAddr)
		return nil
	}

	return fmt.Errorf("Haven't implemented store8 into address 0x%08x with val 0x%02x", addr, val)
}
