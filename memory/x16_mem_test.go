package memory

import "testing"

func TestAllRAMLocations2048K(t *testing.T) {
	mem := NewX16Memory(X2048K)
	var block uint16
	var count uint16
	var largeCount uint32 = 0xA000

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

			if mem.ToLargeMemory().LoadLarge(largeCount) != b {
				t.Fatalf("RAM Block %d not written correctly (2048K, large address)", b)
			}

			if mem.GetStatistics(count) != 3 {
				t.Fatalf("RAM Block %d statistic not written correctly (2048K)", b)
			}

			largeCount++
		}
	}
}

func TestAllRAMLocations512K(t *testing.T) {
	mem := NewX16Memory(X512K)
	var block uint8
	var count uint16
	var largeCount uint32 = 0xA000

	for largeCount = 0xA000; largeCount < 0xA000+(512*1024); largeCount++ {
		block := (largeCount - 0xA000) / 8192
		mem.ToLargeMemory().StoreLarge(largeCount, uint8(block))
	}

	largeCount = 0xA000

	for block = 0; block < 64; block++ {
		mem.Store(0, block)

		for count = 0xA000; count < 0xC000; count++ {
			if mem.Load(count) != block {
				t.Fatalf("RAM Block %d not written correctly (512K)", block)
			}

			if mem.ToLargeMemory().LoadLarge(largeCount) != block {
				t.Fatalf("RAM Block %d not written correctly (512K, large address)", block)
			}

			if mem.GetStatistics(count) != 3 {
				t.Fatalf("RAM Block %d statistic not written correctly (512K)", block)
			}

			largeCount++
		}
	}
}

func TestAllROMLocations(t *testing.T) {
	mem := NewX16Memory(X2048K)
	var block uint8
	var count uint32 = 0xC000
	var largeCount uint32 = 0xA000 + (2048 * 1024)

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

			if mem.ToLargeMemory().LoadLarge(largeCount) != block {
				t.Fatalf("ROM Block %d not written correctly (large address)", block)
			}

			if mem.GetStatistics(uint16(count)) != 3 {
				t.Fatalf("ROM Block %d statistic not written correctly", b)
			}

			largeCount++
		}
	}
}
