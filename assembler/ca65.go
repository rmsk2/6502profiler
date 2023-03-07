package assembler

import (
	"fmt"
	"os/exec"
)

func NewCa65(path string, srcDir string, binDir string, testDir string) *SimpleAsmImpl {
	return &SimpleAsmImpl{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
		parseLine:    parseOneLineCa65,
		genCmd:       makeCa65Cmd,
	}
}

func parseOneLineCa65(line string) (uint16, string, error) {
	return 0, "", fmt.Errorf("ca65 does not provide a possibility to create a label file")
}

func makeCa65Cmd(asmBin, sourceDir string, outName string, progName string) *exec.Cmd {
	return exec.Command(asmBin, "-Wa", fmt.Sprintf("-I,%s", sourceDir), "-o", outName, "-C", "c64-asm.cfg", "--start-addr", "0x0800", progName)
}
