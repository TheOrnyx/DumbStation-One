package memory

// basically the same as the range used in psx-guide but here :P

type Range struct {
	start  uint32
	length uint32
}

var (
	BIOS_RANGE    = Range{start: 0xbfc00000, length: 512 * 1024}
	MEM_CONTROL   = Range{start: 0x1f801000, length: 36}
	RAM_SIZE      = Range{start: 0x1f801060, length: 4} // guide says to ignore
	CACHE_CONTROL = Range{start: 0xfffe0130, length: 4} // the cache control
	RAM_RANGE = Range{start: 0xa0000000, length: 2*1024*1024}
)

// Contains whether or not addr is inside range
func (r *Range) Contains(addr uint32) (uint32, bool) {
	if addr >= r.start && addr < r.start+r.length {
		return addr - r.start, true
	} else {
		return 0xfeed, false // yummy
	}

}
