package assembler

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

type Tass64 struct {
	binPath      string
	srcDir       string
	binDir       string
	testDir      string
	errorMessage string
}

func NewTass64(path string, srcDir string, binDir string, testDir string) *Tass64 {
	return &Tass64{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
	}
}

func (t *Tass64) parseOneLine(line string) (uint16, string, error) {
	r := regexp.MustCompile(`^\s*([[:word:]]+)\s+= [$]([[:xdigit:]]{1,4})(\s.*)?$`)

	matches := r.FindStringSubmatch(line)

	if matches == nil {
		return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
	}

	// Can not fail as the regex ensures that only valid hex numbers are parsed
	res, _ := strconv.ParseUint(matches[2], 16, 16)

	return uint16(res), matches[1], nil
}

func (t *Tass64) ParseLabelFile(fileName string) (map[uint16][]string, error) {
	return ParseLabelFile(fileName, t.parseOneLine)
}

func (t *Tass64) GetErrorMessage() string {
	return t.errorMessage
}

func (t *Tass64) Assemble(fileName string) (string, error) {
	mlProg := path.Join(t.binDir, fmt.Sprintf("%s.bin", fileName))
	mlSrc := path.Join(t.testDir, fileName)
	cmd := exec.Command(t.binPath, "-I", t.srcDir, "-o", mlProg, "-a", mlSrc)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.errorMessage = string(out)
		return "", fmt.Errorf("unable to assemble '%s'", fileName)
	}

	return mlProg, nil
}
