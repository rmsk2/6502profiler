package memory

import "testing"

func TestAllRAMLocations2048K(t *testing.T) {
	mem := NewX16Memory(X2048K)
	var block uint16
	var count uint16

	for block = 0; block < 256; block++ {
		b := uint8(block)
		mem.Store(0, b)

		for count = 0xA000; count < 0xC000; count++ {
			mem.Store(count, b)
		}
	}

	for block = 0; block < 256; block++ {
		b := uint8(block)
		mem.Store(0, b)

		for count = 0xA000; count < 0xC000; count++ {
			if mem.Load(count) != b {
				t.Fatalf("RAM Block %d not written correctly (2048K)", b)
			}

			if mem.GetStatistics(count) != 2 {
				t.Fatalf("RAM Block %d statistic not written correctly (2048K)", b)
			}
		}
	}
}

func TestAllRAMLocations512K(t *testing.T) {
	mem := NewX16Memory(X512K)
	var block uint8
	var count uint16

	for block = 0; block < 64; block++ {
		mem.Store(0, block)

		for count = 0xA000; count < 0xC000; count++ {
			mem.Store(count, uint8(block))
		}
	}

	for block = 0; block < 64; block++ {
		mem.Store(0, block)

		for count = 0xA000; count < 0xC000; count++ {
			if mem.Load(count) != block {
				t.Fatalf("RAM Block %d not written correctly (512K)", block)
			}

			if mem.GetStatistics(count) != 2 {
				t.Fatalf("RAM Block %d statistic not written correctly (2048K)", block)
			}
		}
	}
}

func TestAllROMLocations(t *testing.T) {
	mem := NewX16Memory(X2048K)
	var block uint8
	var count uint32 = 0xC000

	for block = 0; block < 32; block++ {
		mem.Store(1, block)
		for count = 0xC000; count < 0x10000; count++ {
			mem.Store(uint16(count), block)
		}
	}

	for block = 0; block < 32; block++ {
		b := uint8(block)
		mem.Store(1, b)
		for count = 0xC000; count < 0x10000; count++ {
			if mem.Load(uint16(count)) != block {
				t.Fatalf("ROM Block %d not written correctly", b)
			}

			if mem.GetStatistics(uint16(count)) != 2 {
				t.Fatalf("ROM Block %d statistic not written correctly", b)
			}
		}
	}
}
