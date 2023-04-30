package assembler

import (
	"testing"
)

func TestLineParsingTass(t *testing.T) {
	s1 := "XN_SQUARE	= $294f"
	s2 := "DEFAULT_INIT_REAL= $2963"
	s3 := "RES_X		= 320"

	addr, label, err := parseOneLineTass(s1)
	if (addr != 0x294f) || (label != "XN_SQUARE") || (err != nil) {
		t.Fatalf("Matching first test line failed: %d, '%s'", addr, label)
	}

	addr, label, err = parseOneLineTass(s2)
	if (addr != 0x2963) || (label != "DEFAULT_INIT_REAL") || (err != nil) {
		t.Fatalf("Matching second test line failed: %d, '%s'", addr, label)
	}

	addr, label, err = parseOneLineTass(s3)
	if (addr != 320) || (label != "RES_X") || (err != nil) {
		t.Fatalf("Matching third test line failed %d, '%s'", addr, label)
	}
}
