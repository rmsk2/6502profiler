package util

var petsciiTable map[byte]byte = map[byte]byte{}

func init() {
	var i byte

	for i = 32; i <= 90; i++ {
		petsciiTable[i] = i
	}

	for i = 193; i <= 218; i++ {
		petsciiTable[i] = i - 96
	}

	petsciiTable[0x0A] = 0x0A
	petsciiTable[0x0D] = 0x0D
}

func PetsciiToAscii(t byte) byte {
	res, ok := petsciiTable[t]
	if !ok {
		return 63
	}

	return res
}
