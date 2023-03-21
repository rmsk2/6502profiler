package memory

type F256RevBMemory struct {
	mmuMemCtrl *uint8
	mmuIoCtrl  *uint8

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

func NewF56JrMemory() *F256RevBMemory {
	sysMem := make([]byte, 0x140000)

	res := F256RevBMemory{
		mmuMemCtrl: &sysMem[0],
		mmuIoCtrl:  &sysMem[1],

		systemMemory: sysMem,
		ioMemory:     make([]byte, 32768),
		mLut:         make([]byte, 4*8),

		systemMemorySnap: make([]byte, 0x140000),
		ioMemorySnap:     make([]byte, 32768),
		mLutSnap:         make([]byte, 4*8),

		accessSystem: make([]uint64, 0x140000),
		accessIo:     make([]uint64, 32768),
		accessMLut:   make([]uint64, 4*8),
	}

	return &res
}

func (f *F256RevBMemory) calcIndex(addr uint16) (*uint8, *uint64) {
	return nil, nil
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
}

func (f *F256RevBMemory) TakeSnaphot() {
	copy(f.systemMemorySnap, f.systemMemory)
	copy(f.ioMemorySnap, f.ioMemory)
	copy(f.mLutSnap, f.mLut)
}

func (f *F256RevBMemory) RestoreSnapshot() {
	copy(f.systemMemory, f.systemMemorySnap)
	copy(f.ioMemory, f.ioMemorySnap)
	copy(f.mLut, f.mLutSnap)
}

func (f *F256RevBMemory) Close() {

}
