package commands

import (
	"6502profiler/cpu"
	"6502profiler/emuconfig"
	"6502profiler/memory"
	"6502profiler/profiler"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const strategyMedian = "median"

func determineCutOffCalc(strategy *string, p float64) profiler.CutOffCalc {
	var ctOff profiler.CutOffCalc

	if *strategy != strategyMedian {
		ctOff = func(m memory.Memory, start, end uint16) uint64 {
			return profiler.CutOffAbsoluteValue(m, start, end, p)
		}
	} else {
		ctOff = func(m memory.Memory, start, end uint16) uint64 {
			return profiler.CutOffMedian(m, start, end, p)
		}
	}

	return ctOff
}

func parseDumpParams(param string) (uint16, uint16, error) {
	if param == "" {
		return 0, 0, nil
	}

	r := regexp.MustCompile(`^([0-9]+):([0-9]+)$`)
	matches := r.FindStringSubmatch(param)
	if matches == nil {
		return 0, 0, fmt.Errorf("unable to parse dump parameters")
	}

	dumpAddress, err := strconv.ParseUint(matches[1], 10, 16)
	if err != nil {
		return 0, 0, err
	}

	dumpLen, err := strconv.ParseUint(matches[2], 10, 16)
	if err != nil {
		return 0, 0, err
	}

	if dumpLen == 0 {
		return 0, 0, fmt.Errorf("length has to be at least 1")
	}

	dumpAddress16 := uint16(dumpAddress)
	dumpLen16 := uint16(dumpLen)

	if (dumpAddress16 + dumpLen16 - 1) < dumpAddress16 {
		return 0, 0, fmt.Errorf("dump parameters result in end address being lower than start address")
	}

	return dumpAddress16, dumpLen16, nil
}

func DumpMemory(param string, cpu *cpu.CPU6502) error {
	if param == "" {
		return nil
	}

	dumpAddress, dumpLen, err := parseDumpParams(param)
	if err != nil {
		return err
	}

	memory.Dump(cpu.Mem, dumpAddress, dumpAddress+dumpLen-1)

	return nil
}

func ProfileCommand(arguments []string) error {
	var p float64
	var labels map[uint16][]string
	var err error = nil
	var config *emuconfig.Config = emuconfig.DefaultConfig()
	var processor *cpu.CPU6502

	profileFlags := flag.NewFlagSet("6502profiler profile", flag.ContinueOnError)
	binaryFileName := profileFlags.String("prg", "", "Path to the program to run")
	labelFileName := profileFlags.String("label", "", "Path to the label file generated by the ACME assembler")
	outputFileName := profileFlags.String("out", "", "Path to the out file that holds the generated data")
	percentageCutOff := profileFlags.Uint("prcnt", 10, "Percentage used to determine cut off value")
	strategy := profileFlags.String("strategy", strategyMedian, "Strategy to determine cutoff value")
	configName := profileFlags.String("c", "", "Config file name")
	dumpFlag := profileFlags.String("dump", "", "Dump memory after program has stopped. Format 'startaddr:len'")

	if err = profileFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = emuconfig.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	assembler := config.GetAssembler()

	processor, err = config.NewCpu()
	if err != nil {
		return fmt.Errorf("error processing config: %v", err)
	}
	defer func() { processor.Mem.Close() }()

	statisticRequested := (*outputFileName != "")

	if *binaryFileName == "" {
		return fmt.Errorf("no program specified")
	}

	if _, _, err := parseDumpParams(*dumpFlag); err != nil {
		return err
	}

	if statisticRequested {
		if *labelFileName == "" {
			return fmt.Errorf("a label file has to be specified")
		}

		if *percentageCutOff > 100 {
			return fmt.Errorf("%d is not a valid value for cutoff percentage", *percentageCutOff)
		}

		labels, err = assembler.ParseLabelFile(*labelFileName)
		if err != nil {
			return fmt.Errorf("a problem occurred: %v", err)
		}

		p = float64(*percentageCutOff) / 100.0
	}

	loadAddress, progLen, err := processor.LoadAndRun(*binaryFileName)
	if err != nil {
		return fmt.Errorf("a problem occurred: %v", err)
	}

	fmt.Printf("Program ran for %d clock cycles\n", processor.NumCycles())

	if statisticRequested {
		var ctOff = determineCutOffCalc(strategy, p)

		if err = profiler.DumpStatistics(processor.Mem, *outputFileName, labels, loadAddress, (loadAddress + progLen - 1), ctOff); err != nil {
			return fmt.Errorf("problem generating output file: %v", err)
		}
	}

	err = DumpMemory(*dumpFlag, processor)
	if err != nil {
		return err
	}

	return nil
}
