// utility functions for use in the other parts of the emulator
package utils

// BytesToUint32 Convert 4 bytes to one little endian uint32 value
func BytesToUint32(b0, b1, b2, b3 byte) uint32 {
	return uint32(b0) | (uint32(b1) << 8) | (uint32(b2) << 16) | (uint32(b3) << 24)
}

// Uint32ToBytes convert a uint32 word to it's 4 individual bytes in
// little endian form
func Uint32ToBytes(data uint32) (b0, b1, b2, b3 byte) {
	b0 = byte(data & 0xFF)
	b1 = byte((data >> 8) & 0xFF)
	b2 = byte((data >> 16) & 0xFF)
	b3 = byte((data >> 24) & 0xFF)
	return
}

