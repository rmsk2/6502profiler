package profiler

import (
	"6502profiler/memory"
	"fmt"
	"os"
)

type CutOffCalc func(m memory.Memory, start uint16, end uint16) uint64

func DumpStatistics(m memory.Memory, fileName string, acmeLabels map[uint16][]string, start uint16, end uint16, determineCutOffValue CutOffCalc) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() { f.Close() }()

	cutOff := determineCutOffValue(m, start, end)

	for count := start; count <= end; count++ {
		labels, ok := acmeLabels[count]
		if ok {
			for _, j := range labels {
				fmt.Fprintln(f, j)
			}
		}

		numAccess := m.GetStatistics(count)
		if numAccess != 0 {
			numAccess -= 1
		}

		prefix := "     "
		if numAccess >= cutOff {
			prefix = "###  "
		}

		fmt.Fprintf(f, "%s%04x: %02X %d\n", prefix, count, m.Load(count), numAccess)
	}

	return nil
}
