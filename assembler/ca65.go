package assembler

import (
	"fmt"
	"os/exec"
	"path"
)

type Ca65AsmImpl struct {
	binPath      string
	srcDir       string
	binDir       string
	testDir      string
	errorMessage string
	startAddress uint16
}

func NewCa65(path string, srcDir string, binDir string, testDir string, startAddress uint16) *Ca65AsmImpl {
	return &Ca65AsmImpl{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
		startAddress: startAddress,
	}
}

func (c *Ca65AsmImpl) ParseLabelFile(fileName string) (map[uint16][]string, error) {
	return nil, fmt.Errorf("unsupported: ca65 is unable to create an easily parseable symbol list file")
}

func (c *Ca65AsmImpl) GetErrorMessage() string {
	return c.errorMessage
}

func (c *Ca65AsmImpl) GetDefaultSrc() string {
	return `
.PC02

main:
    brk
`
}

func (c *Ca65AsmImpl) Assemble(fileName string) (string, error) {
	mlProg := path.Join(c.binDir, fmt.Sprintf("%s.bin", fileName))
	mlObj := path.Join(c.binDir, fmt.Sprintf("%s.obj", fileName))
	mlSrc := path.Join(c.testDir, fileName)
	asmCommand := path.Join(c.binPath, "ca65")
	linkCommand := path.Join(c.binPath, "cl65")

	asmCmd := exec.Command(asmCommand,
		"-I", c.srcDir,
		"-o", mlObj,
		mlSrc,
	)

	loadAddress := fmt.Sprintf("0x%04x", c.startAddress)

	linkCmd := exec.Command(linkCommand,
		"-C", "c64-asm.cfg",
		"--start-addr", loadAddress,
		"-o", mlProg,
		mlObj,
	)

	out, err := asmCmd.CombinedOutput()
	if err != nil {
		c.errorMessage = string(out)
		return "", fmt.Errorf("unable to assemble '%s'", fileName)
	}

	out, err = linkCmd.CombinedOutput()
	if err != nil {
		c.errorMessage = string(out)
		return "", fmt.Errorf("unable to link '%s'", fileName)
	}

	return mlProg, nil
}
