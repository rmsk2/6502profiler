package assembler

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

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

func (a *ACME) parseOneLine(line string) (uint16, string, error) {
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
	return ParseLabelFile(fileName, a.parseOneLine)
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
