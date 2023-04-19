package memory

import "testing"

func TestMul161(t *testing.T) {
	mem := NewLinearMemory(65536)
	memWrap := NewMemWrapper(mem, 0xDE00)
	coproc := NewUnsignedCoproc(mem, 0xDE00)

	coproc.RegisterUmul(memWrap)

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

	if memWrap.Load(0xDE07) != 0x25 {
		t.Fatal("Byte 3 is wrong")
	}
}

func TestDiv161(t *testing.T) {
	mem := NewLinearMemory(65536)
	memWrap := NewMemWrapper(mem, 0xDE00)
	coproc := NewUnsignedCoproc(mem, 0xDE00)

	coproc.RegisterUdiv(memWrap)

	memWrap.Store(0xDE08, 0x23)
	memWrap.Store(0xDE09, 0x45)
	memWrap.Store(0xDE0A, 0x67)
	memWrap.Store(0xDE0B, 0x89)

	if memWrap.Load(0xDE14) != 0x01 {
		t.Fatal("Byte 0 is wrong")
	}

	if memWrap.Load(0xDE15) != 0x00 {
		t.Fatal("Byte 1 is wrong")
	}

	if memWrap.Load(0xDE16) != 0x44 {
		t.Fatal("Byte 2 is wrong")
	}

	if memWrap.Load(0xDE17) != 0x44 {
		t.Fatal("Byte 3 is wrong")
	}
}
