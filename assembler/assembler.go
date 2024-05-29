package assembler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Assembler interface {
	Assemble(fileName string) (string, error)
	ParseLabelFile(fileName string) (map[uint16][]string, error)
	GetErrorMessage() string
	GetDefaultSrc() string
}

type LineParseFunc func(string) (uint16, string, error)
type GenCommandFunc func(asmBin string, sourceDir string, outName string, progName string, binDir string, obFile string) *exec.Cmd

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

type SimpleAsmImpl struct {
	binPath      string
	srcDir       string
	binDir       string
	testDir      string
	errorMessage string
	parseLine    LineParseFunc
	genCmd       GenCommandFunc
	defaultProg  string
}

func (s *SimpleAsmImpl) ParseLabelFile(fileName string) (map[uint16][]string, error) {
	return ParseLabelFile(fileName, s.parseLine)
}

func (s *SimpleAsmImpl) GetErrorMessage() string {
	return s.errorMessage
}

func (s *SimpleAsmImpl) GetDefaultSrc() string {
	return s.defaultProg
}

func (s *SimpleAsmImpl) Assemble(fileName string) (string, error) {
	mlProg := path.Join(s.binDir, fmt.Sprintf("%s.bin", fileName))
	mlObj := path.Join(s.binDir, fmt.Sprintf("%s.obj", fileName))
	mlSrc := path.Join(s.testDir, fileName)
	cmd := s.genCmd(s.binPath, s.srcDir, mlProg, mlSrc, s.binDir, mlObj)

	out, err := cmd.CombinedOutput()
	if err != nil {
		s.errorMessage = string(out)
		return "", fmt.Errorf("unable to assemble '%s'", fileName)
	}

	return mlProg, nil
}
