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
		case 0: // base
			return channel.base
		case 4: // Block control register
			return channel.BlockControl()
		case 8:
			return channel.Control()
		default:
			log.Panicf("unhandled DMA read at: 0x%08x, minor:0x%04x", offset, minor)
		}

	case 7: // Common DMA registers
		switch minor {
		case 0:
			return b.dma.Control()
		case 4:
			return b.dma.Interrupt()
		default:
			log.Panicf("unhandled DMA read at: 0x%08x, minor:0x%04x", offset, minor)
		}
	}

	log.Panicf("unhandled DMA read at: 0x%08x", offset)
	return 0x00 // shouldn't be reached
}

// SetDMAReg DMA register write
func (b *Bus) SetDMAReg(offset, val uint32) {
	major := (offset & 0x70) >> 4
	minor := offset & 0xf
	var activePort Port
	portFound := false

	switch major {
	// per channel registers
	case 0,1,2,3,4,5,6: // TODO - this is gross, kill later cbf rn
		port := PortFromIndex(major)
		channel := b.dma.GetChannelRef(port)
		switch minor {
		case 0: // base
			channel.SetBase(val)
		case 4: // block control
			channel.SetBlockControl(val)
		case 8:
			channel.SetControl(val)
		default:
			log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x, minor:0x%04x", val, offset, minor)
		}

		if channel.IsActive() {
			activePort = port
			portFound = true
		}

	case 7: // Common DMA registers
		switch minor {
		case 0:
			b.dma.SetControl(val)
		case 4:
			b.dma.SetInterrupt(val)
		default:
			log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x, minor:0x%04x", val, offset, minor)
		}

	default:
		log.Panicf("Unhandled DMA write: 0x%08x into 0x%08x", val, offset)
	}

	if portFound {
		b.doDMA(activePort)
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
		case 4: // gpustat
			// set bit 26, 27 and 28 to signal that gpu is ready for
			// DMA and CPU access
			return 0x1c000000, nil
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


//////////////////////////////////
// Perform DMA Transfer methods //
//////////////////////////////////

// doDMA execute a DMA transfer to a port
func (b *Bus) doDMA(port Port)  {
	channel := b.dma.GetChannelRef(port)

	// Don't care about splitting stuff up, do everything in one pass
	switch channel.syncMode {
	case linkedListMode: // not implemented
		log.Panic("(Not implemented yet) Linked list DMA mode not supported")
	default:
		b.doDMABlock(port)
	}
}

// doDMABlock do dma block transfer
func (b *Bus) doDMABlock(port Port)  {
	channel := b.dma.GetChannelRef(port)
	increment := channel.Step()

	addr := channel.base

	var remainingSize uint32 // transfer size in words
	if size, notLinked := channel.TransferSize(); notLinked {
		remainingSize = size
	} else { // this shouldn't happen
		log.Panic("Couldn't figure out DMA block transfer size")
	}

	for remainingSize > 0 {
		currentAddr := addr & 0x1ffffc

		switch channel.transferDir {
		case dirFromRam:
			log.Panicf("Unhandled DMA direction")

		case dirToRam:
			srcWord := b.getDMASrcWord(port, addr, remainingSize)

			b.ram.store32(currentAddr, srcWord)
		}

		addr = addr + uint32(increment)
		remainingSize -= 1
	}

	channel.Done()
}

// getDMASrcWord get the source word for DMA transfer at that point
// This is a seperate method cuz I didn't really wanna have like 3
// more switch statements in doDMABlock
func (b *Bus) getDMASrcWord(port Port, addr, remainingSize uint32) uint32 {
	var srcWord uint32
	switch port {
	case PortOtc: // clear ordering table
		if remainingSize == 1 {
			srcWord = 0xffffff
		} else {
			srcWord = (addr + 4) & 0x1fffff
		}

	default:
		log.Panicf("Unhandled DMA source port: %v", port)
	}

	return srcWord
}
