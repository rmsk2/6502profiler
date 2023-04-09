package memory

import "testing"

func TestEditMlut(t *testing.T) {
	mem := NewF56JrMemory(false)

	mem.Store(0, 0)
	// Set MLUT 0
	mem.SetMlut(0, []byte{0, 1, 2, 3, 4, 5, 6, 7})

	// Set target address of MLUT editing to defined value
	mem.mLut[3*lutSize+1] = 0x11

	// Here byte 9 is "normal" memory
	mem.Store(9, 0x42)

	// Turn on edit mode for MLUT 3
	mem.Store(0, 0b10110000)
	// Here byte 9 is mapped to second byte of MLUT 3
	mem.Store(9, 0xFF)
	// Turn MLUT editing off and activate MLUT 0
	mem.Store(0, 0b00000000)

	// Check whether MLUT 3 was modified as expected
	if mem.mLut[3*lutSize+1] != 0xFF {
		t.Fatal("Modifying MLUT did not work")
	}

	// Check whether normal memory has retained defined value
	if mem.Load(9) != 0x42 {
		t.Fatal("Address decoding does not work")
	}
}

func TestDefaultMemory(t *testing.T) {
	mem := NewF56JrMemory(true)
	// MLUT 1 is active
	mem.Store(0, 1)
	// Set MLUT 1
	mem.SetMlut(1, []byte{16, 17, 18, 19, 20, 21, 22, 23})

	var memAddress uint = 128*1024 + 4193
	mem.systemMemory[memAddress] = 0x39

	if mem.Load(4193) != 0x39 {
		t.Fatal("Read wrong value via MLUT")
	}

	mem.Store(4193, 0xBD)

	if mem.systemMemory[memAddress] != 0xBD {
		t.Fatal("Read wrong value via direct access")
	}

	// Change MLUT 1
	mem.SetMlut(1, []byte{17, 16, 18, 19, 20, 21, 22, 23})

	// memAddress 4193 is now one 8K bank "further up"
	memAddress = 128*1024 + uint(bankSize) + 4193

	mem.systemMemory[memAddress] = 0x45

	if mem.Load(4193) != 0x45 {
		t.Fatal("Read wrong value via MLUT on second try")
	}

	mem.Store(4193, 0xEF)

	if mem.systemMemory[memAddress] != 0xEF {
		t.Fatal("Read wrong value via direct access on second try")
	}
}

func TestIoAccess(t *testing.T) {
	mem := NewF56JrMemory(true)
	var memAddress uint = 128*1024 + 8
	// Bank nr.6 has value 16
	mem.SetMlut(2, []byte{22, 17, 18, 19, 20, 21, 16, 23})

	// Write defined value into IO bank 1
	mem.ioMemory[1*bankSize+8] = 0x23
	// MLUT 2 is active
	mem.Store(0, 2)
	// Memory IO page 1 is active
	mem.Store(1, 1)

	if mem.Load(0xC008) != 0x23 {
		t.Fatal("Read wrong value from memory IO area")
	}

	mem.Store(0xC008, 0x54)

	if mem.ioMemory[1*bankSize+8] != 0x54 {
		t.Fatal("Writing to IO memory area failed")
	}

	// Disble IO Memory, activate entry nr. 6 of active MLUT
	mem.Store(1, 0b00000101)
	// Write to system memory
	mem.systemMemory[memAddress] = 0xAB

	// This should now read from system memory and not from IO memory
	if mem.Load(0xC008) != 0xAB {
		t.Fatal("Disabling IO memory did not work")
	}
}

func TestEditActiveMLUT(t *testing.T) {
	mem := NewF56JrMemory(true)
	mem.SetMlut(1, []byte{16, 17, 18, 19, 20, 21, 22, 23})
	// MLUT 1 is active and is editable
	mem.Store(0, 145)

	mem.Store(0x2000+723, 17)

	mem.Store(9, 25)
	mem.Store(0x2000+723, 44)

	mem.Store(9, 17)

	if mem.Load(0x2000+723) != 17 {
		t.Fatal("Should have read 17")
	}

	mem.Store(9, 25)
	if mem.Load(0x2000+723) != 44 {
		t.Fatal("Should have read 44")
	}
}
