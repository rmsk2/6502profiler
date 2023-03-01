package memory

import (
	"6502profiler/util"
	"fmt"
	"regexp"
)

type ToAsciiFunc func(byte) byte

type PrinterProcessor struct {
	conv ToAsciiFunc
}

// NewPrinterProcessorFromConfig parses a config string of the form "printer:petscii" and
// creates the corresponding PrinterProcessor struct. The string after the colon specifies
// the encoding.
func NewPrinterProcessorFromConfig(configEntry string) (MemWrapper, bool) {
	r := regexp.MustCompile(`^printer:([0-9a-zA-Z]+)$`)

	matches := r.FindStringSubmatch(configEntry)
	if matches == nil {
		return nil, false
	}

	var res MemWrapper = nil

	if matches[1] == "petscii" {
		return NewPetsciiPrinter(), true
	}

	return res, false
}

func NewPetsciiPrinter() *PrinterProcessor {
	return &PrinterProcessor{
		conv: util.PetsciiToAscii,
	}
}

func (p *PrinterProcessor) Write(b uint8) {
	fmt.Printf("%c", p.conv(b))
}

func (p *PrinterProcessor) Close() {
}
