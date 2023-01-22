package cpu

import "testing"

// lda #$42
// brk
func TestLDAImmediate(t *testing.T) {
	prog := []byte{0xA9, 0x42, 0x00}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		if c.A != 0x42 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	})

	if res == false {
		t.Fatal("LDA immediate does not work")
	}
	if err != nil {
		t.Fatalf("LDA immediate does not work: %v", err)
	}
}

// lda #00
// brk
func TestLDAImmediate0(t *testing.T) {
	prog := []byte{0xA9, 0x00, 0x00}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		if c.A != 0x00 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	})

	if res == false {
		t.Fatal("LDA immediate does not work")
	}
	if err != nil {
		t.Fatalf("LDA immediate does not work: %v", err)
	}
}

// lda #$81
// brk
func TestLDAImmediateNeg(t *testing.T) {
	prog := []byte{0xA9, 0x81, 0x00}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		if c.A != 0x81 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		return true
	})

	if res == false {
		t.Fatal("LDA immediate does not work")
	}
	if err != nil {
		t.Fatalf("LDA immediate does not work: %v", err)
	}
}

// lda $0804
// brk
// !byte 0x72
// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsolute(t *testing.T) {
	prog := []byte{0xAD, 0x04, 0x08, 0x00, 0x72}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		return c.A == 0x72
	})

	if res == false {
		t.Fatal("LDA absolute does not work")
	}
	if err != nil {
		t.Fatalf("LDA absolute does not work: %v", err)
	}
}

// ldx #6
// lda $0800, x
// brk
// !byte 0x72
// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsoluteX(t *testing.T) {
	prog := []byte{0xA2, 0x06, 0xBD, 0x00, 0x08, 0x00, 0x72}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		return c.A == 0x72
	})

	if res == false {
		t.Fatal("LDA absolute X does not work")
	}
	if err != nil {
		t.Fatalf("LDA absolute X does not work: %v", err)
	}
}

// ldy #6
// lda $0800, y
// brk
// !byte 0x72
// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsoluteY(t *testing.T) {
	prog := []byte{0xA0, 0x06, 0xB9, 0x00, 0x08, 0x00, 0x72}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		return c.A == 0x72
	})

	if res == false {
		t.Fatal("LDA absolute Y does not work")
	}
	if err != nil {
		t.Fatalf("LDA absolute Y does not work: %v", err)
	}
}

// lda #<$0800
// sta $12
// lda #>$0800
// sta $13
// ldy #$0d
// lda ($12),y
// brk
// DATA
// !byte $72
func TestLDAIndirectIdxY(t *testing.T) {
	prog := []byte{0xa9, 0x00, 0x85, 0x12, 0xa9, 0x08, 0x85, 0x13, 0xa0, 0x0d, 0xb1, 0x12, 0x00, 0x72}
	res, err := testSingleInstruction(Model6502, prog, func(c *CPU6502) bool {
		return c.A == 0x72
	})

	if res == false {
		t.Fatal("LDA indirect index Y does not work")
	}
	if err != nil {
		t.Fatalf("LDA indirect index Y does not work: %v", err)
	}
}
