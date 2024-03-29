package memory

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type StdOutProcessor struct {
	lineLength uint
	charCount  uint
}

// NewStdOutProcessorFromConfig parses a config string of the form "stdout:16" and
// creates the corresponding StdOutProcessor struct. The number after the colon
// defines the line length og the output.
func NewStdOutProcessorFromConfig(configEntry string) (MemWrapper, bool) {
	r := regexp.MustCompile(`^stdout:([0-9]+)$`)

	matches := r.FindStringSubmatch(configEntry)
	if matches == nil {
		return nil, false
	}

	lineLength, _ := strconv.ParseUint(matches[1], 10, 32)

	res := NewStdOutProcessor(uint(lineLength))

	return res, true
}

func NewStdOutProcessor(lineLen uint) *StdOutProcessor {
	return &StdOutProcessor{
		lineLength: lineLen,
		charCount:  0,
	}
}

func (s *StdOutProcessor) Write(b uint8) {
	if (s.charCount != 0) && ((s.charCount % s.lineLength) == 0) {
		fmt.Println()
	}
	fmt.Printf("%02X ", b)
	s.charCount++
}

func (s *StdOutProcessor) SetBaseMem(m Memory) {

}

// -------------------------------

type StdOutBinaryProcessor struct {
}

// NewStdOutBinaryProcessorFromConfig parses a config string of the form "stdout:bin" and
// creates the corresponding StdOutBinaryProcessor struct.
func NewStdOutBinaryProcessorFromConfig(configEntry string) (MemWrapper, bool) {
	if configEntry != "stdout:bin" {
		return nil, false
	}

	return NewStdOutBinaryProcessor(), true
}

func NewStdOutBinaryProcessor() *StdOutBinaryProcessor {
	return &StdOutBinaryProcessor{}
}

func (s *StdOutBinaryProcessor) Write(b uint8) {
	outB := []byte{b}
	_, err := os.Stdout.Write(outB)
	if err != nil {
		panic(err)
	}
}

func (s *StdOutBinaryProcessor) SetBaseMem(m Memory) {

}
