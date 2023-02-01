package memory

import (
	"fmt"
	"os"
)

type FileProcessor struct {
	f *os.File
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

func (p *FileProcessor) Close() error {
	return p.f.Close()
}
