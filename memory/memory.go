package memory

import (
	"fmt"
	"strconv"
	"strings"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, b uint8)
	GetStatistics(address uint16) uint64
	ClearStatistics()
	TakeSnaphot()
	RestoreSnapshot()
	Close()
}

func Dump(m Memory, start uint16, end uint16) {
	byteCount := 0
	crlfWritten := false
	printableChars := ""

	lineDone := func() {
		crlfWritten = true
		byteCount = 0
		fmt.Println(printableChars + "|")
	}

	for i := start; i <= end; i++ {
		crlfWritten = false

		if byteCount == 0 {
			fmt.Printf("$%04x  ", i)
			printableChars = " |"
		}

		if byteCount == 8 {
			fmt.Printf(" ")
		}

		data := m.Load(i)

		if (data < 128) && (strconv.IsPrint(rune(data))) {
			printableChars += string([]byte{data})
		} else {
			printableChars += "."
		}

		fmt.Printf("%02x ", data)

		if byteCount == 15 {
			lineDone()
		} else {
			byteCount++
		}
	}

	if !crlfWritten {
		maxSpaces := strings.Repeat(" ", 46)

		plusOneSpace := 0
		if byteCount <= 8 {
			plusOneSpace++
		}

		fmt.Print(maxSpaces[:(16-byteCount)*3+plusOneSpace])
		lineDone()
	}

	fmt.Printf("$%04x\n", end)
}
