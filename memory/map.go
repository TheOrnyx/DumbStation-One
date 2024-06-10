package memory

// basically the same as the range used in psx-guide but here :P

type Range struct {
	start  uint32
	length uint32
}

var (
	RAM_RANGE     = Range{start: 0x00000000, length: 2 * 1024 * 1024}
	BIOS_RANGE    = Range{start: 0x1fc00000, length: 512 * 1024}
	SYS_CONTROL   = Range{start: 0x1f801000, length: 36}
	RAM_SIZE      = Range{start: 0x1f801060, length: 4} // guide says to ignore
	CACHE_CONTROL = Range{start: 0xfffe0130, length: 4} // the cache control
	SPU_RANGE     = Range{start: 0x1f801c00, length: 640}
	EXPANSION_1   = Range{start: 0x1f000000, length: 8192*1024} // TODO - check i have no idea
	EXPANSION_2   = Range{start: 0x1f802000, length: 66}
	IRQ_CONTROL   = Range{start: 0x1f801070, length: 8} // interrupr request
	TIMERS_RANGE  = Range{start: 0x1f801100, length: 48} // TODO - check, idk the fucking memory map is confusing as shit
	DMA_RANGE     = Range{start: 0x1f801080, length: 0x80}
)

// Contains whether or not addr is inside range
func (r *Range) Contains(addr uint32) (uint32, bool) {
	if addr >= r.start && addr < r.start+r.length {
		return addr - r.start, true
	} else {
		return 0xfeed, false // yummy
	}

}

// Mask array to strip region bits from the address.  The mask is
// seleced using the 3 MSBs of the address so each entry matches the
// 512kB of the address space. KSEG2 is not touched
var REGION_MASK [8]uint32 = [8]uint32{
	// KUSEG: 2048MB
	0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff,
	// KSEG0: 512MB
	0x7fffffff,
	// KSEG1: 512MB
	0x1fffffff,
	// KSEG2: 1024MB
	0xffffffff, 0xffffffff}

// MaskRegion mask a CPU address to remove the region bits
func MaskRegion(addr uint32) uint32 {
	index := addr >> 29

	return addr & REGION_MASK[index]
}
