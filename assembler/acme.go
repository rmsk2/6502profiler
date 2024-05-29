package assembler

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const defaultDrvAcme string = `
* = $0800
!cpu 65c02

main
    brk
`

func NewACME(path string, srcDir string, binDir string, testDir string) *SimpleAsmImpl {
	return &SimpleAsmImpl{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
		parseLine:    parseOneLineAcme,
		genCmd:       makeAcmeCmd,
		defaultProg:  defaultDrvAcme,
	}
}

func parseOneLineAcme(line string) (uint16, string, error) {
	r := regexp.MustCompile(`^\s+([[:word:]]+)\s+= [$]([[:xdigit:]]{1,4})(\s.*)?$`)

	matches := r.FindStringSubmatch(line)

	if matches == nil {
		return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
	}

	// Can not fail as the regex ensures that only valid hex numbers are parsed
	res, _ := strconv.ParseUint(matches[2], 16, 16)

	return uint16(res), matches[1], nil
}

func makeAcmeCmd(asmBin, sourceDir string, outName string, progName string, binDir string, obFile string) *exec.Cmd {
	return exec.Command(asmBin, "-I", sourceDir, "-o", outName, "-f", "cbm", progName)
}
