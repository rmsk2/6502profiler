package memory

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, b uint8)
	GetStatistics(address uint16) uint64
}

type CutOffCalc func(m Memory, start uint16, end uint16) uint64

func Dump(m Memory, start uint16, end uint16) {
	temp := make([]byte, end-start+1)
	var index uint16 = 0

	for count := start; count <= end; count++ {
		temp[index] = m.Load(count)
		index++
	}

	fmt.Print(hex.Dump(temp))
}

func DumpStatistics(m Memory, fileName string, acmeLabels map[uint16][]string, start uint16, end uint16, determineCutOffValue CutOffCalc) error {
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
