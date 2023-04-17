package memory

import (
	"regexp"
	"strconv"
	"strings"
)

type F256UnsignedMultiplierWrapper struct {
	inAddrStart  uint16
	outAddrStart uint16
	offset       uint16
	mem          Memory
}

func NewUMultiplierFromConfig(conf string) (MemWrapper, bool) {
	components := strings.Split(conf, ":")

	if len(components) != 2 {
		return nil, false
	}

	if components[0] != "F256UMul" {
		return nil, false
	}

	r := regexp.MustCompile(`^([0-9]{1,5})#([0-9]{1,5})#([0-9])$`)

	matches := r.FindStringSubmatch(components[1])
	if matches == nil {
		return nil, false
	}

	outAddr, _ := strconv.ParseUint(matches[2], 10, 16)
	if outAddr >= 65536 {
		return nil, false
	}

	inAddr, _ := strconv.ParseUint(matches[1], 10, 16)
	if inAddr >= 65536 {
		return nil, false
	}

	offset, _ := strconv.ParseUint(matches[3], 10, 16)
	if offset >= 4 {
		return nil, false
	}

	return NewUMultiplier(uint16(inAddr), uint16(outAddr), uint16(offset)), true
}

func NewUMultiplier(inAddr uint16, outAddr uint16, o uint16) *F256UnsignedMultiplierWrapper {
	return &F256UnsignedMultiplierWrapper{
		inAddrStart:  inAddr,
		outAddrStart: outAddr,
		offset:       o,
	}
}

func (f *F256UnsignedMultiplierWrapper) Close() {

}

func (f *F256UnsignedMultiplierWrapper) Write(data uint8) {
	f.mem.Store(f.inAddrStart+f.offset, data)
	oper1 := uint(f.mem.Load(f.inAddrStart+1))*256 + uint(f.mem.Load(f.inAddrStart))
	oper2 := uint(f.mem.Load(f.inAddrStart+3))*256 + uint(f.mem.Load(f.inAddrStart+2))

	res := oper1 * oper2
	var i uint16

	for i = 0; i < 4; i++ {
		f.mem.Store(f.outAddrStart+i, uint8(res&0xFF))
		res >>= 8
	}
}

func (f *F256UnsignedMultiplierWrapper) SetBaseMem(m Memory) {
	f.mem = m
}
