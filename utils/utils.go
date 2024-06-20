// utility functions for use in the other parts of the emulator
package utils

import (
	"math"
)

// BytesToUint32 Convert 4 bytes to one little endian uint32 value
func BytesToUint32(b0, b1, b2, b3 byte) uint32 {
	return uint32(b0) | (uint32(b1) << 8) | (uint32(b2) << 16) | (uint32(b3) << 24)
}

// BytesToUint16 convert 2 bytes to one little endian uint16 value
func BytesToUint16(b0, b1 byte) uint16 {
	return uint16(b0) | (uint16(b1) << 8)
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

// AddSigned16 add two signed 16 bit values together and return the result and
// whether or not they caused an overflow
//
// FIXME - this is really likely wrong, like genuinely I don't trust this at all
func AddSigned16(val1, val2 uint32) (result uint32, overflowed bool) {
	sVal1 := int32(val1)
	sVal2 := int32(val2)

	res := int64(sVal1) + int64(sVal2)
	if res > math.MaxInt32 || res < math.MinInt32 {
		return 0, true
	}

	return uint32(int32(res)), false

	// origBit := (signedVal >> 31)

	// res := signedVal + signedImm
	// if res >> 31 != origBit {
	// 	return 0x00, true
	// }

	// return uint32(res), false
}

// SubSigned16 subtract two signed 16-bit values and return the result
// and whether they overflowed
//
// FIXME - check this cuz it probs ain't right
func SubSigned16(val1, val2 uint32) (result uint32, overflowed bool) {
	sVal1 := int32(val1)
	sVal2 := int32(val2)

	res := int64(sVal1) - int64(sVal2)
	if res > math.MaxInt32 || res < math.MinInt32 {
		return 0, true
	}

	return uint32(int32(res)), false
}

// // BoolToInt[T constrains.Integer] convert conditional to 0 or 1 of type
// func BoolToInt[T constraints.Integer](cond bool) T {
// 	res := 0

// 	if cond {
// 		res = 1
// 	}

// 	return T(res)
// }

// BoolToUint32 convert cond to 0 or 1 and return
//
// Apparently this is the fastest way to do it according to
// https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3 but i
// should probs check later
func BoolToUint32(cond bool) uint32 {
	var r uint32
	if cond {
		r = 1
	} else {
		r = 0
	}

	return r
}
