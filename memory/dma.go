package memory

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/utils"
)

// direct memory access
// TODO - document this better based on the docs - guide is confusing
type Dma struct {
	control uint32 // DMA control register

	// DICR register at 0x1F8010F4h
	irqControl uint8 // bits [0:6] of interrupts registers control channel 0-6 completion? NOTE - guide says unknown and that this is 0-5 but specs says otherwise so idk
	// Bits [7:14] are unused
	forceIRQ      bool  // when set interrupt is active unconditionally (even if irqEnable is false) bit 15?
	chanIRQEnable uint8 // IRQ enable for each channel - bits[16:22]
	enableIRQ     bool  // master Interrupt request enable bit 23
	chanIRQFlags  uint8 // IRQ flags for each channel - bits[24:30]

	// The 7 channel instances
	channels [7]ChannelControl
}

type Port int

const (
	PortMdecIn  Port = 0 // macroblock decoder input
	PortMdecOut Port = 1 // Macroblock decoder input
	PortGpu     Port = 2 // GPU port
	PortCdRom   Port = 3 // CD-ROM driver
	PortSpu     Port = 4 // Sound Processing Unit
	PortPio     Port = 5 // Extension Port
	PortOtc     Port = 6 // Used to clear ordering table
)

// PortFromIndex return a port based on the given index
func PortFromIndex(index uint32) Port {
	switch index {
	case 0: return PortMdecIn
	case 1: return PortMdecOut
	case 2: return PortGpu
	case 3: return PortCdRom
	case 4: return PortSpu
	case 5: return PortPio
	case 6: return PortOtc
	default:
		log.Panicf("Invalid port index: %v", index)
	}

	return 7 // shouldn't happen
}

// NewDMA create and return a new dma
func NewDMA() Dma {
	return Dma{control: 0x07654321, channels: [7]ChannelControl{}}
}

// Control return the control register
func (d *Dma) Control() uint32 {
	return d.control
}

// SetControl set control to val
func (d *Dma) SetControl(val uint32) {
	d.control = val
}

// GetChannelRef return a reference to an entry in the channels array using a port num
func (d *Dma) GetChannelRef(port Port) *ChannelControl {
	return &d.channels[port]
}

///////////////////////////
// IRQ methods and stuff //
///////////////////////////

// IRQ return the status of the DMA interrupt basically does what bit 31 would be
func (d *Dma) IRQ() bool {
	channelIRQ := d.chanIRQFlags & d.chanIRQEnable

	return d.forceIRQ || (d.enableIRQ && channelIRQ != 0)
}

// IRQU return status of DMA interrupt as a uint32 value
func (d *Dma) IRQU() uint32 {
	return utils.BoolToUint32(d.IRQ())
}

// Interrupt get the uint32 value of the DICR DMA interrupt register
func (d *Dma) Interrupt() uint32 {
	var res uint32

	res |= uint32(d.irqControl)
	res |= d.forceIRQU() << 15
	res |= uint32(d.chanIRQEnable) << 16
	res |= d.enableIRQU() << 23
	res |= uint32(d.chanIRQFlags) << 24
	res |= d.IRQU() << 31

	return res
}

// SetInterrupt set value of interrupt register
func (d *Dma) SetInterrupt(val uint32) {
	d.irqControl = uint8(val & 0x7f) // masked bits 0-6

	d.forceIRQ = (val>>15)&1 != 0

	d.chanIRQEnable = uint8((val >> 16) & 0x7f)

	d.enableIRQ = (val>>23)&1 != 0

	// writing 1 to flag resets it
	ack := uint8((val >> 24) & 0x3f)
	d.chanIRQFlags &= (^ack)
}

// enableIRQU return the enableIRQ bool as uint32
func (d *Dma) enableIRQU() uint32 {
	return utils.BoolToUint32(d.enableIRQ)
}

// forceIRQU return forceIRQ as uint32
func (d *Dma) forceIRQU() uint32 {
	return utils.BoolToUint32(d.forceIRQ)
}

///////////////////////////////
// DMA Channel Control stuff //
///////////////////////////////

// D_CHCR DMA Channel Control (R/W)
//
// TODO - maybe simplify this, it kinda sucks, like maybe just do a uint32 and have each thing be a method idk
type ChannelControl struct {
	transferDir   uint8 // Transfer direction (0 = device to RAM, 1 = RAM to device) (Bit 0)
	stepIncrement uint8 // MADR increment per step (0 = +4, 1 = -4) (bit 1)
	// Bits 2-7 are unused
	chopping bool  // bit 8 - chopping enabled - NOTE the specs say it can be different so check later
	syncMode uint8 // bits 9-10 Transfer mode (SyncMode) 0=burst, 1=slice, 2=linked-list, 3=reserved
	// Bits 11-15 are unused
	chopDMASize uint8 // Chopping DMA windows size (1 << N words) - (Bits 16-18)
	// Bit 19 unused
	chopCPUSize uint8 // Chopping CPU windows size (1 << N cycles) - (Bits 20-22)
	// Bit 23 unused
	enabled bool // Start transfer (false=stopped/completed, 1=start/busy) - (bit 24)
	// bits 25-27 unsused
	forceStart bool // force transfer start without waiting for DREQ - (Bit 28)
	// rest aren't important so will just represent them with a uint8
	upper uint8 // the unimportant bits 29-31

	// DMA start address
	base uint32

	// Block stuff
	blockSize uint16 // size of a block in words
	blockCount uint16 // block count, used only when 'syncMode' is 'sliceMode' (guide says request mode)
}

// syncMode constants
const (
	burstMode      = 0
	sliceMode      = 1
	linkedListMode = 2
	reservedMode   = 3
)

// newChannelControl create and return a new channelControl object
func newChannelControl() ChannelControl {
	return ChannelControl{
		transferDir:   0, // toram
		stepIncrement: 0, // +4
		chopping:      false,
		syncMode:      burstMode, // manual sync mode
		chopDMASize:   0,
		chopCPUSize:   0,
		enabled:       false,
		forceStart:    false,
	}
}

// Control return channel control register as uint32
//
// FIXME - this is gross and probably rlly unoptimized - change later
func (c *ChannelControl) Control() uint32 {
	var r uint32

	r |= uint32(c.transferDir) << 0
	r |= uint32(c.stepIncrement) << 1
	r |= utils.BoolToUint32(c.chopping) << 8
	r |= uint32(c.syncMode) << 9
	r |= uint32(c.chopDMASize) << 16
	r |= uint32(c.chopCPUSize) << 20
	r |= utils.BoolToUint32(c.enabled) << 24
	r |= utils.BoolToUint32(c.forceStart) << 28
	r |= uint32(c.upper) << 29

	return r
}

// SetControl set the control registers fields
func (c *ChannelControl) SetControl(val uint32) {
	c.transferDir = uint8(val & 0x01)
	c.stepIncrement = uint8((val >> 1) & 0x01)

	c.chopping = (val>>8)&0x01 != 0

	switch (val >> 9) & 0x03 {
	case 0:
		c.syncMode = burstMode
	case 1:
		c.syncMode = sliceMode
	case 2:
		c.syncMode = linkedListMode
	default: // TODO - check if reserved is supposed to be used?
		log.Panicf("Unknown syncMode 0x%04x", (val>>9)&0x03)
	}

	c.chopDMASize = uint8((val >> 16) & 0x07)
	c.chopCPUSize = uint8((val >> 20) & 0x07)

	c.enabled = (val>>24)&0x01 != 0
	c.forceStart = (val>>24)&0x01 != 0
	c.upper = uint8((val >> 29) & 0x03)
}

// SetBase set channel base, only bits [0:23] are significant so only
// 16mb are addressable by the DMA
func (c *ChannelControl) SetBase(val uint32)  {
	c.base = val & 0xffffff
}

// BlockControl get the value of the blockcontrol register
func (c *ChannelControl) BlockControl() uint32 {
	bs := uint32(c.blockSize)
	bc := uint32(c.blockCount)

	return (bc << 16) | bs
}

// SetBlockControl set the value of the block control register
func (c *ChannelControl) SetBlockControl(val uint32)  {
	c.blockSize = uint16(val)
	c.blockCount = uint16(val >> 16)
}
