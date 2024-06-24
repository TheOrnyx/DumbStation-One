/*
 * The CDROM package for handling all of the parsing and handling of CD data etc
 */
package cdrom

import (
	"fmt"
	"os"

	"github.com/TheOrnyx/psx-go/log"
)

type CDROM struct {
	Data         []byte // The executable data (TODO - make this into an actual type interface alter)
	status       Status // Index/Status register (0x1f801800)
	intFlagReg   uint8  // Interrupt flag register
	intEnableReg uint8  // Interrupt enable register
}

type Status uint8 // The Index/Status Register - TODO - maybe convert to struct

// Status read and return the status
func (c *CDROM) Status() uint8 {
	return c.status.index()
}

// writeStatus write to the status register
func (c *CDROM) writeStatus(val uint8)  {
	c.status = Status(val & 0x03)
}

// Index get the index from stat register
func (c *CDROM) Index() uint8 {
	return c.status.index()
}

// index get the index from status register (bits 0-1)
func (s Status) index() uint8 {
	return uint8(s) & 3
}

// NewCDROM Create and return a new CDROM
func NewCDROM(path string) (CDROM, error) {
	var cd CDROM
	data, err := os.ReadFile(path)
	if err != nil {
		return cd, fmt.Errorf("Failed to open file: %v", err)
	}

	cd.Data = data
	return cd, nil
}

// ReadResponse read from the Response FIFO - TODO complete
func (c *CDROM) ReadResponse() uint8 {
	log.Warn("(Not implemented yet) attempted read from REsponse FIFO")
	return 0x00
}

// ReadData Read data from the CDROM Data FIFO
func (c *CDROM) ReadData() uint8 {
	log.Warn("(Not implemented yet) attempted read from Data FIFO")
	return 0x00
}

// LoadByte read byte in CDROM register at addr
func (c *CDROM) LoadByte(offset uint32) uint8 {
	switch offset {
	case 0:
		return c.Status()
	case 1:
		return c.ReadResponse()
	case 2:
		return c.ReadData()
	case 3:
		if c.Index() & 1 == 1 { // 1 or 3
			return 0xe0 | c.intFlagReg
		}

		return c.intEnableReg
	}

	log.Panicf("Unhandled CDROM Read from offset %v", offset)
	return 0xF
}

// StoreByte store given byte in CDROM register
func (c *CDROM) StoreByte(offset uint32, val uint8) {
	switch (c.Index() << 2) | uint8(offset) {
	// index 0
	case 0:
		c.writeStatus(val)
	case 1:
		c.WriteCMD(val)
	case 2:
		c.writeParam(val)
	case 3:
		c.writeRequest(val)
	// Index 1
	case 4:
		c.writeStatus(val)
	case 5: // ignore this i think?
		log.Infof("Ignoring write to CDROM Sound Map Data Out") 
	case 6:
		c.writeIntEnable(val)
	case 7:
		c.writeIntFlag(val)
	// Index 2
	case 8:
		c.writeStatus(val)
	case 9: // ignore as well
		log.Infof("Ignoring write to CDROM Sound Map Coding Info")
	case 10:
		c.writeLeftToLeftVol(val)
	case 11:
		c.writeLeftToRightVol(val)
	// Index 3
	case 12:
		c.writeStatus(val)
	case 13:
		c.writeRightToRightVol(val)
	case 14:
		c.writeRightToLeftVol(val)
	case 15:
		c.writeAudioVolApply(val)
	default:
		log.Panicf("Unhandled CDROM Store into offset %v with val 0x%02x", offset, val)
	}
}

// WriteCMD write to command register (should run a command)
func (c *CDROM) WriteCMD(val uint8)  {
	log.Panicf("(Not implemented yet) attempted write to command reg with val 0x%02x", val)
}

// writeParam write to the parameter FIFO
func (c *CDROM) writeParam(val uint8)  {
	log.Warnf("(Not implemented yet) attempted write to Parameter FIFO with val 0x%02x", val)
}

// writeRequest write to the request register
func (c *CDROM) writeRequest(val uint8)  {
	log.Warnf("(Not implemented yet) attempted write to Request Register with val 0x%02x", val)
}

// writeIntEnable write to the interrupt enable register
func (c *CDROM) writeIntEnable(val uint8)  {
	c.intEnableReg = val
}

// writeIntFlag write to interrupt flag register
func (c *CDROM) writeIntFlag(val uint8)  {
	c.intFlagReg = val
}

// writeLeftToLeftVol Audio Volume for Left-CD-Out to Left-SPU-Input
func (c *CDROM) writeLeftToLeftVol(val uint8)  {
	log.Infof("(Not implemented yet) write Audio Volume for Left-CD-Out to Left-SPU-Input with val %d", val)
}

// writeLeftToRightVol Audio Volume for Left-CD-Out to Right-SPU-Input
func (c *CDROM) writeLeftToRightVol(val uint8)  {
	log.Infof("(Not implemented yet) write Audio Volume for Left-CD-Out to Right-SPU-Input with val %d", val)
}

// writeRightToRightVol Audio Volume for Right-CD-Out to Right-SPU-Input
func (c *CDROM) writeRightToRightVol(val uint8)  {
	log.Infof("(Not implemented yet) write Audio Volume for Right-CD-Out to Right-SPU-Input with val %d", val)
}

// writeRightToLeftVol Audio Volume for Right-CD-Out to Left-SPU-Input
func (c *CDROM) writeRightToLeftVol(val uint8)  {
	log.Infof("(Not implemented yet) write Audio Volume for Right-CD-Out to Left-SPU-Input with val %d", val)
}

// writeAudioVolApply Apply changes to volume (By writing bit5=1)
func (c *CDROM) writeAudioVolApply(val uint8)  {
	log.Infof("(Not implemented yet) write Audio Volume Apply Changes with val %d", val)
}
