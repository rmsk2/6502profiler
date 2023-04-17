package memory

import "testing"

func TestMul161(t *testing.T) {
	mem := NewLinearMemory(65536)
	memWrap := NewMemWrapper(mem, 0xDE)

	memWrap.AddWrapper(0xDE00, NewUMultiplier(0xDE00, 0xDE04, 0))
	memWrap.AddWrapper(0xDE01, NewUMultiplier(0xDE00, 0xDE04, 1))
	memWrap.AddWrapper(0xDE02, NewUMultiplier(0xDE00, 0xDE04, 2))
	memWrap.AddWrapper(0xDE03, NewUMultiplier(0xDE00, 0xDE04, 3))

	memWrap.Store(0xDE00, 0x23)
	memWrap.Store(0xDE01, 0x45)
	memWrap.Store(0xDE02, 0x67)
	memWrap.Store(0xDE03, 0x89)

	if memWrap.Load(0xDE04) != 0x15 {
		t.Fatal("Byte 0 is wrong")
	}

	if memWrap.Load(0xDE05) != 0x8c {
		t.Fatal("Byte 1 is wrong")
	}

	if memWrap.Load(0xDE06) != 0x1b {
		t.Fatal("Byte 2 is wrong")
	}

	if memWrap.Load(0xD007) != 0x25 {
		t.Fatal("Byte 3 is wrong")
	}
}
