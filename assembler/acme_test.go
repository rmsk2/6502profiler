package assembler

import (
	"testing"
)

func TestLineParsing(t *testing.T) {
	s1 := "	resetParametersHires	= $df0	; ?"
	s2 := "	LOOKUP_SCRATCH3	= $fd	; unused"
	s3 := "	STEP_X	= $d95"

	addr, label, err := parseOneLine(s1)
	if (addr != 0x0df0) || (label != "resetParametersHires") || (err != nil) {
		t.Fatalf("Matching first test line failed: %d, '%s'", addr, label)
	}

	addr, label, err = parseOneLine(s2)
	if (addr != 0x00fd) || (label != "LOOKUP_SCRATCH3") || (err != nil) {
		t.Fatalf("Matching second test line failed: %d, '%s'", addr, label)
	}

	addr, label, err = parseOneLine(s3)
	if (addr != 0x0d95) || (label != "STEP_X") || (err != nil) {
		t.Fatalf("Matching third test line failed %d, '%s'", addr, label)
	}
}
