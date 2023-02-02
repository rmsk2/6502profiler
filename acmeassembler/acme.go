package acmeassembler

import (
	"bufio"
	"fmt"
	"os"
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

func ParseLabelFile(fileName string) (map[uint16][]string, error) {
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
