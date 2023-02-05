package profiler

import (
	"6502profiler/memory"
	"sort"
)

func CutOffAbsoluteValue(m memory.Memory, start uint16, end uint16, p float64) uint64 {
	temp := map[uint64]bool{}

	for count := start; count <= end; count++ {
		temp[m.GetStatistics(count)] = true
	}

	keys := []uint64{}

	for i := range temp {
		keys = append(keys, i)
	}

	lenKeys := len(keys)
	if lenKeys == 0 {
		return 0
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	l := float64(lenKeys)
	cutOffIndex := int(l * (1.0 - p))

	return keys[cutOffIndex]
}

func CutOffMedian(m memory.Memory, start uint16, end uint16, p float64) uint64 {
	// length of memory area is at least 1

	temp := []uint64{}

	for count := start; count <= end; count++ {
		temp = append(temp, m.GetStatistics(count))
	}

	sort.Slice(temp, func(i, j int) bool { return temp[i] < temp[j] })

	l := float64(len(temp))
	cutOffIndex := int(l * (1.0 - p))

	return temp[cutOffIndex]
}
