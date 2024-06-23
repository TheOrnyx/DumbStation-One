package memory

type AccessWidth uint32
const (
	Byte AccessWidth = 1
	HalfWord AccessWidth = 2
	Word AccessWidth = 4
)

type Addressable interface {
	width() AccessWidth // Retrieve width of the access
	fromUint32(val uint32) Addressable // Build addressable value from uint32.
	asUint32() uint32 // Retrieve value of Addressable as uint32
}

type AddressableUint8 uint8

func (a AddressableUint8) width() AccessWidth {
	return Byte
}

func (a AddressableUint8) fromUint32(val uint32) Addressable {
	return AddressableUint8(val)
}

func (a AddressableUint8) asUint32() uint32 {
	return uint32(a)
}

type AddressableUint16 uint16

func (a AddressableUint16) width() AccessWidth {
	return HalfWord
}

func (a AddressableUint16) fromUint32(val uint32) Addressable {
	return AddressableUint16(val)
}

func (a AddressableUint16) asUint32() uint32 {
	return uint32(a)
}

type AddressableUint32 uint32

func (a AddressableUint32) width() AccessWidth {
	return Word
}

func (a AddressableUint32) fromUint32(val uint32) Addressable {
	return a
}

func (a AddressableUint32) asUint32() uint32 {
	return uint32(a)
}
