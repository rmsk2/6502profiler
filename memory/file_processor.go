package memory

import (
	"fmt"
	"os"
	"strings"
)

type FileProcessor struct {
	f *os.File
}

// NewFileProcFomConfig parses a config string of the form "file:apfel2.bin" and
// creates the corresponding FileProcessor struct.
func NewFileProcFromConfig(conf string) (MemWrapper, bool) {
	components := strings.Split(conf, ":")

	if len(components) != 2 {
		return nil, false
	}

	if components[0] != "file" {
		return nil, false
	}

	f, err := NewFileProcessor(components[1])
	return f, err == nil
}

func NewFileProcessor(fileName string) (*FileProcessor, error) {
	ftemp, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return &FileProcessor{
		f: ftemp,
	}, nil
}

func (p *FileProcessor) Write(b uint8) {
	_, err := p.f.Write([]byte{b})
	if err != nil {
		panic(fmt.Sprintf("Can not write file: %v", err))
	}
}

func (p *FileProcessor) Close() {
	err := p.f.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing data file: %v", err))
	}
}
