package assembler

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const defaultDrv64Tass string = `
* = $0800
.cpu "w65c02"

main
    brk
`

func NewTass64(path string, srcDir string, binDir string, testDir string) *SimpleAsmImpl {
	return &SimpleAsmImpl{
		binPath:      path,
		srcDir:       srcDir,
		binDir:       binDir,
		testDir:      testDir,
		errorMessage: "",
		parseLine:    parseOneLineTass,
		genCmd:       makeTassCmd,
		defaultProg:  defaultDrv64Tass,
	}
}

func parseOneLineTass(line string) (uint16, string, error) {
	var res uint64
	var err error
	hexOrNot := regexp.MustCompile(`^\s*[[:word:]]+\s*= [$].*$`)
	matches := hexOrNot.FindStringSubmatch(line)

	if matches != nil {
		// A dollar sign appears to the right of the equal sign => we look for a hex number
		rHex := regexp.MustCompile(`^\s*([[:word:]]+)\s*= [$]([[:xdigit:]]{1,4})(\s.*)?$`)

		matches = rHex.FindStringSubmatch(line)
		if matches == nil {
			return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
		}

		// Can not fail as the regex ensures that only valid hex numbers are parsed
		res, _ = strconv.ParseUint(matches[2], 16, 16)
	} else {
		// No dollar sign appears to the right of the equal sign => we look for a decimal number
		rDec := regexp.MustCompile(`^\s*([[:word:]]+)\s*= ([[:digit:]]{1,5})(\s.*)?$`)

		matches = rDec.FindStringSubmatch(line)
		if matches == nil {
			return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
		}

		// Can fail as the regex can not ensure that any five digit number fits into a 16 bit int
		res, err = strconv.ParseUint(matches[2], 10, 16)
		if err != nil {
			return 0, "", fmt.Errorf("can not parse label file line '%s'", line)
		}
	}

	return uint16(res), matches[1], nil
}

func makeTassCmd(asmBin, sourceDir string, outName string, progName string, binDir string, obFile string) *exec.Cmd {
	return exec.Command(asmBin, "-I", sourceDir, "-o", outName, "-a", progName)
}
