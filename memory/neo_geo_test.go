package memory

import "testing"

var testTrack uint16 = NeoGeoRegisterPage + 0xFE
var testSector uint16 = NeoGeoRegisterPage + 0xFF

func TestMaskBits(t *testing.T) {
	if checkSectorBits(9) != 8 {
		t.Fatal("Wrong number of sector bits (9)")
	}

	if checkSectorBits(0) != 8 {
		t.Fatal("Wrong number of sector bits (0)")
	}

	if checkSectorBits(5) != 5 {
		t.Fatal("Wrong number of sector bits (5)")
	}

	if checkSectorBits(8) != 8 {
		t.Fatal("Wrong number of sector bits (8)")
	}

	if checkSectorBits(7) != 7 {
		t.Fatal("Wrong number of sector bits (8)")
	}

	if calcSectorMask(8) != 0xFF {
		t.Fatal("Wrong number sector mask (8)")
	}

	if calcSectorMask(5) != 0x1F {
		t.Fatal("Wrong number sector mask (5)")
	}

	if calcSectorMask(7) != 0x7F {
		t.Fatal("Wrong number sector mask (5)")
	}
}

func TestNeoGeo1(t *testing.T) {
	var neo Memory = NewNeoGeo(testTrack, 5)

	// testSector == 0, testTrack == 0
	neo.Store(NeoGeoRamPage+1, 42)
	// testSector == 0, testTrack == 1
	neo.Store(testTrack, 1)
	neo.Store(NeoGeoRamPage+1, 43)

	if neo.Load(NeoGeoRamPage+1) != 43 {
		t.Fatal("Expected to load 43 from NeoGeoRamPage+1")
	}

	neo.Store(testTrack, 0)
	if neo.Load(NeoGeoRamPage+1) != 42 {
		t.Fatal("Expected to load 42 from NeoGeoRamPage+1")
	}

	// testTrack == 0, testSector == 1
	neo.Store(testSector, 1)
	neo.Store(NeoGeoRamPage+1, 44)
	// testTrack == 1, testSector == 1
	neo.Store(testTrack, 1)
	neo.Store(NeoGeoRamPage+1, 45)

	neo.Store(testTrack, 0)
	if neo.Load(NeoGeoRamPage+1) != 44 {
		t.Fatal("Expected to load 44 from NeoGeoRamPage+1")
	}

	neo.Store(testTrack, 1)
	if neo.Load(NeoGeoRamPage+1) != 45 {
		t.Fatal("Expected to load 45 from NeoGeoRamPage+1")
	}

	neo.Store(testTrack, 0)
	neo.Store(testSector, 0)
	if neo.Load(NeoGeoRamPage+1) != 42 {
		t.Fatal("Expected to load 42 from NeoGeoRamPage+1 again")
	}

	neo.Store(testTrack, 64)
	if neo.Load(NeoGeoRamPage+1) != 42 {
		t.Fatal("Expected to load 42 from NeoGeoRamPage+1 again and again")
	}

	neo.Store(testSector, 32)
	if neo.Load(NeoGeoRamPage+1) != 42 {
		t.Fatal("Expected to load 42 from NeoGeoRamPage+1 again and again and again")
	}
}

func TestCalcAddress(t *testing.T) {
	neo := NewNeoGeo(testTrack, 5)
	neo.Store(testTrack, 63)
	neo.Store(testSector, 31)
	neo.Store(NeoGeoRamPage+0xFF, 42)

	if neo.calcIndexRaw(NeoGeoRamPage+0xFF) != 524287 {
		t.Fatal("Expected maximal address to be 524287")
	}

	if neo.calcIndexRaw(NeoGeoRamPage+0xF0) != (524287 - 15) {
		t.Fatal("Expected address to be 524287 - 15")
	}

	neo = NewNeoGeo(testTrack, 8)
	neo.Store(testTrack, 63)
	neo.Store(testSector, 0xFF)
	neo.Store(NeoGeoRamPage+0xFF, 42)

	if neo.calcIndexRaw(NeoGeoRamPage+0xFF) != 4194303 {
		t.Fatal("Expected maximal address to be 4194304")
	}

	if neo.calcIndexRaw(NeoGeoRamPage+0xF0) != (4194303 - 15) {
		t.Fatal("Expected address to be 4194304 - 15")
	}

	neo = NewNeoGeo(testTrack, 7)
	neo.Store(testTrack, 63)
	neo.Store(testSector, 0x7F)
	neo.Store(NeoGeoRamPage+0xFF, 42)

	if neo.calcIndexRaw(NeoGeoRamPage+0xFF) != 2097151 {
		t.Fatal("Expected maximal address to be 2097151")
	}

	if neo.calcIndexRaw(NeoGeoRamPage+0xF0) != (2097151 - 15) {
		t.Fatal("Expected address to be 2097151 - 15")
	}
}
