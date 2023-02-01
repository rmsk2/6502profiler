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

func Dump(m Memory, start uint16, end uint16) {
	temp := make([]byte, end-start+1)
	var index uint16 = 0

	for count := start; count <= end; count++ {
		temp[index] = m.Load(count)
		index++
	}

	fmt.Print(hex.Dump(temp))
}

func DumpStatistics(m Memory, fileName string, start uint16, end uint16) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() { f.Close() }()

	for count := start; count <= end; count++ {
		fmt.Fprintf(f, "%04x: %02X %d\n", count, m.Load(count), m.GetStatistics(count)-1)
	}

	return nil
}
