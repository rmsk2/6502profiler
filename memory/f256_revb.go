package memory

type F256RevBMemory struct {
	mmuMemCtrl uint8
	mmuIoCtrl  uint8

	mmuMemCtrlSnap uint8
	mmuIoCtrlSnap  uint8

	accessMmuMemCtrl uint64
	accessMmuIoCtrl  uint64

	systemMemory []byte
	ioMemory     []byte
	mLut         []byte

	systemMemorySnap []byte
	ioMemorySnap     []byte
	mLutSnap         []byte

	accessSystem []uint64
	accessIo     []uint64
	accessMLut   []uint64
}

const bankSize uint16 = 8192
const lutSize = 8
const numLuts = 4
const numIoBanks = 4

func NewF56JrMemory(expandedMem bool) *F256RevBMemory {
	memSize := 0x100000

	if expandedMem {
		memSize = 0x140000
	}

	res := F256RevBMemory{
		mmuMemCtrl:   0,
		mmuIoCtrl:    0,
		systemMemory: make([]byte, memSize),
		ioMemory:     make([]byte, numIoBanks*bankSize),
		mLut:         make([]byte, numLuts*lutSize),

		mmuMemCtrlSnap:   0,
		mmuIoCtrlSnap:    0,
		systemMemorySnap: make([]byte, memSize),
		ioMemorySnap:     make([]byte, numIoBanks*bankSize),
		mLutSnap:         make([]byte, numLuts*lutSize),

		accessMmuMemCtrl: 0,
		accessMmuIoCtrl:  0,
		accessSystem:     make([]uint64, memSize),
		accessIo:         make([]uint64, numIoBanks*bankSize),
		accessMLut:       make([]uint64, numLuts*lutSize),
	}

	res.SetMlut(0, []byte{0, 1, 2, 3, 4, 5, 6, 7})
	res.SetMlut(1, []byte{0, 1, 2, 3, 4, 5, 6, 7})
	res.SetMlut(2, []byte{0, 1, 2, 3, 4, 5, 6, 7})
	res.SetMlut(3, []byte{0, 1, 2, 3, 4, 5, 6, 7})

	// MLUT 0 is active
	// IO bank 0 is active

	return &res
}

func (f *F256RevBMemory) calcDefault(activeLut []byte, loBits uint16, hiBits uint16) uint {
	lower := uint(loBits)
	var upper uint = (uint(activeLut[hiBits]) << numLoBits)

	return upper | lower
}

const loBitMask uint16 = 0b0001111111111111
const numLoBits uint16 = 13
const lutEditFlagMask uint8 = 0b10000000
const editLutMask uint8 = 0b00110000
const activeLutMask uint8 = 0b00000011
const ioDisableMask uint8 = 0b00000100
const activeIoBankMask uint8 = 0b00000011

func (f *F256RevBMemory) calcIndex(addr uint16) (*uint8, *uint64) {
	// lower 13 bits
	loBits := addr & loBitMask
	hiBits := addr >> numLoBits
	activeLutNum := f.mmuMemCtrl & activeLutMask

	activeLut := f.mLut[activeLutNum*lutSize : (activeLutNum+1)*lutSize]

	switch {
	case addr == 0:
		return &f.mmuMemCtrl, &f.accessMmuMemCtrl
	case addr == 1:
		return &f.mmuIoCtrl, &f.accessMmuIoCtrl
	case (addr >= 8) && (addr <= (8 + lutSize - 1)):
		if (f.mmuMemCtrl & lutEditFlagMask) != 0 {
			editLut := uint16((f.mmuMemCtrl & editLutMask) >> 4)
			return &f.mLut[editLut*lutSize+(addr-8)], &f.accessMLut[editLut*lutSize+(addr-8)]
		} else {
			idx := f.calcDefault(activeLut, loBits, hiBits)
			return &f.systemMemory[idx], &f.accessSystem[idx]
		}
	case (addr >= 0xC000) && (addr <= 0xDFFF):
		if (f.mmuIoCtrl & ioDisableMask) == 0 {
			ioBank := uint16(f.mmuIoCtrl & activeIoBankMask)
			idx := ioBank*bankSize + addr - 0xC000
			return &f.ioMemory[idx], &f.accessIo[idx]
		} else {
			idx := f.calcDefault(activeLut, loBits, hiBits)
			return &f.systemMemory[idx], &f.accessSystem[idx]
		}
	default:
		idx := f.calcDefault(activeLut, loBits, hiBits)
		return &f.systemMemory[idx], &f.accessSystem[idx]
	}
}

func (f *F256RevBMemory) Load(address uint16) uint8 {
	memPtr, statPtr := f.calcIndex(address)
	(*statPtr)++
	return *memPtr
}

func (f *F256RevBMemory) Store(address uint16, b uint8) {
	memPtr, statPtr := f.calcIndex(address)
	(*statPtr)++
	*memPtr = b
}

func (f *F256RevBMemory) GetStatistics(address uint16) uint64 {
	_, statPtr := f.calcIndex(address)
	return *statPtr
}

func (f *F256RevBMemory) ClearStatistics() {
	for count := 0; count < len(f.accessSystem); count++ {
		f.accessSystem[count] = 0
	}

	for count := 0; count < len(f.accessIo); count++ {
		f.accessIo[count] = 0
	}

	for count := 0; count < len(f.accessMLut); count++ {
		f.accessMLut[count] = 0
	}

	f.accessMmuMemCtrl = 0
	f.accessMmuIoCtrl = 0
}

func (f *F256RevBMemory) TakeSnapshot() {
	copy(f.systemMemorySnap, f.systemMemory)
	copy(f.ioMemorySnap, f.ioMemory)
	copy(f.mLutSnap, f.mLut)
	f.mmuMemCtrlSnap = f.mmuMemCtrl
	f.mmuIoCtrlSnap = f.mmuIoCtrl
}

func (f *F256RevBMemory) RestoreSnapshot() {
	copy(f.systemMemory, f.systemMemorySnap)
	copy(f.ioMemory, f.ioMemorySnap)
	copy(f.mLut, f.mLutSnap)
	f.mmuMemCtrl = f.mmuMemCtrlSnap
	f.mmuIoCtrl = f.mmuIoCtrlSnap
}

func (f *F256RevBMemory) SetMlut(lutNum uint8, mlutData []byte) {
	lutNum = lutNum & 0b00000011
	copy(f.mLut[lutNum*lutSize:(lutNum+1)*lutSize], mlutData)
}
