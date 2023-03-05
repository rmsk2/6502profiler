package assembler

import (
	"bufio"
	"fmt"
	"os"
)

type Assembler interface {
	Assemble(fileName string) (string, error)
	ParseLabelFile(fileName string) (map[uint16][]string, error)
	GetErrorMessage() string
}

type LineParseFunc func(string) (uint16, string, error)

func ParseLabelFile(fileName string, parseOneLine LineParseFunc) (map[uint16][]string, error) {
	result := make(map[uint16][]string)

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() { f.Close() }()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		addr, label, err := parseOneLine(fileScanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error reading label file: %v", err)
		}

		_, ok := result[addr]
		if ok {
			result[addr] = append(result[addr], label)
		} else {
			result[addr] = []string{label}
		}
	}

	return result, nil
}
