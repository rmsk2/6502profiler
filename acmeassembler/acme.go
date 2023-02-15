package acmeassembler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

func parseOneLine(line string) (uint16, string, error) {
	r := regexp.MustCompile(`^\s+([[:word:]]+)\s+= [$]([[:xdigit:]]{1,4})(\s.*)?$`)

	matches := r.FindStringSubmatch(line)

	if matches == nil {
		return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
	}

	// Can not fail as the regex ensures that only valid hex numbers are parsed
	res, _ := strconv.ParseUint(matches[2], 16, 16)

	return uint16(res), matches[1], nil
}

func (a *ACME) ParseLabelFile(fileName string) (map[uint16][]string, error) {
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

type ACME struct {
	binPath      string
	srcDir       string
	binDir       string
	testDir      string
	errorMessage string
}

func NewACME(path string, srcDir string, binDir string, testDir string) *ACME {
	return &ACME{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
	}
}

func (a *ACME) GetErrorMessage() string {
	return a.errorMessage
}

func (a *ACME) Assemble(fileName string) (string, error) {
	mlProg := path.Join(a.binDir, fmt.Sprintf("%s.bin", fileName))
	mlSrc := path.Join(a.testDir, fileName)
	cmd := exec.Command(a.binPath, "-I", a.srcDir, "-o", mlProg, "-f", "cbm", mlSrc)

	out, err := cmd.CombinedOutput()
	if err != nil {
		a.errorMessage = string(out)
		return "", fmt.Errorf("unable to assemble '%s'", fileName)
	}

	return mlProg, nil
}
