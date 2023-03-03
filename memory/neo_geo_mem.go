package memory

const NeoGeoRamPage = 0xDE00
const NeoGeoRegisterPage = NeoGeoRamPage + 0x100

type NeoGeoRam struct {
	sectorMask uint8
	sectorBits uint

	trackPtr  *byte
	sectorPtr *byte
	baseMem   []byte
	neoGeo    []byte

	baseMemSnapshot []byte
	neoGeoSnapshot  []byte

	statBase   []uint64
	statNeoGeo []uint64
}

func NewNeoGeo(trackPtrAddress uint16, sectorBits uint) *NeoGeoRam {
	var geoSize uint = 1
	geoSize <<= checkSectorBits(sectorBits) + 14
	// Make track address even and mask out MSB
	trackPtrAddress &= 0x00FE
	// Make sure track address is in the range DF00-DFFF
	trackPtrAddress = trackPtrAddress | NeoGeoRegisterPage

	baseRam := make([]byte, 65536)

	res := &NeoGeoRam{
		sectorMask:      calcSectorMask(sectorBits),
		sectorBits:      checkSectorBits(sectorBits),
		trackPtr:        &baseRam[trackPtrAddress],
		sectorPtr:       &baseRam[trackPtrAddress+1],
		baseMem:         baseRam,
		baseMemSnapshot: make([]byte, 65536),

		neoGeo:         make([]byte, geoSize),
		neoGeoSnapshot: make([]byte, geoSize),

		statBase:   make([]uint64, 65536),
		statNeoGeo: make([]uint64, geoSize),
	}

	*res.sectorPtr = 0
	*res.trackPtr = 0

	return res
}

func checkSectorBits(sectorBits uint) uint {
	if (sectorBits == 0) || (sectorBits > 8) {
		sectorBits = 8
	}

	return sectorBits
}

func calcSectorMask(sectorBits uint) uint8 {
	sectorBits = checkSectorBits(sectorBits)
	var res uint8 = 1

	return (res << uint8(sectorBits)) - 1
}

func (n *NeoGeoRam) Close() {

}

func (n *NeoGeoRam) calcIndexRaw(address uint16) uint {
	byteOffset := uint(address & 0xFF)

	geoAddr := uint(*n.trackPtr & 0x3F)
	geoAddr <<= n.sectorBits
	geoAddr |= uint(*n.sectorPtr & n.sectorMask)
	geoAddr <<= 8
	geoAddr |= byteOffset

	return geoAddr
}

func (n *NeoGeoRam) calcIndex(address uint16) (*uint8, *uint64) {
	switch {
	case (address < NeoGeoRamPage) || (address >= NeoGeoRegisterPage):
		return &n.baseMem[address], &n.statBase[address]
	default:
		geoAddr := n.calcIndexRaw(address)
		return &n.neoGeo[geoAddr], &n.statNeoGeo[geoAddr]
	}
}

func (n *NeoGeoRam) TakeSnaphot() {
	copy(n.baseMemSnapshot, n.baseMem)
	copy(n.neoGeoSnapshot, n.neoGeo)
}

func (n *NeoGeoRam) RestoreSnapshot() {
	copy(n.baseMem, n.baseMemSnapshot)
	copy(n.neoGeo, n.neoGeoSnapshot)
}

func (n *NeoGeoRam) ClearStatistics() {
	for i := 0; i < len(n.statBase); i++ {
		n.statBase[i] = 0
	}

	for i := 0; i < len(n.statNeoGeo); i++ {
		n.statNeoGeo[i] = 0
	}
}

func (n *NeoGeoRam) Load(address uint16) uint8 {
	mem, stat := n.calcIndex(address)
	(*stat)++
	return *mem
}

func (n *NeoGeoRam) Store(address uint16, b uint8) {
	mem, stat := n.calcIndex(address)
	(*stat)++
	*mem = b
}

func (n *NeoGeoRam) GetStatistics(address uint16) uint64 {
	_, stat := n.calcIndex(address)
	return *stat
}
