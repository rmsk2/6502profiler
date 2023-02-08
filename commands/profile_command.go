package commands

import (
	"6502profiler/acmeassembler"
	"6502profiler/cpu"
	"6502profiler/memory"
	"6502profiler/profiler"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
)

func ProfileCommand(arguments []string) error {
	const strategyMedian = "median"
	var p float64
	var labels map[uint16][]string
	var err error = nil
	var config *cpu.Config = cpu.EmptyConfig()
	var processor *cpu.CPU6502

	profileFlags := flag.NewFlagSet("6502profiler profile", flag.ContinueOnError)
	binaryFileName := profileFlags.String("prg", "", "Path to the program to run")
	labelFileName := profileFlags.String("label", "", "Path to the label file generated by the ACME assembler")
	outputFileName := profileFlags.String("out", "", "Path to the out file that holds the generated data")
	percentageCutOff := profileFlags.Uint("prcnt", 10, "Percentage used to determine cut off value")
	strategy := profileFlags.String("strategy", strategyMedian, "Strategy to determine cutoff value")
	configName := profileFlags.String("c", "", "Config file name")

	// mem := memory.NewMemWrapper(memory.NewLinearMemory(16384), 0x2D00)
	// picProc := memory.NewPicProcessor(320, 200)
	// mem.AddSpecialWriteAddress(0x2DDD, picProc.SetPoint)
	//mem := memory.NewLinearMemory(65536)

	if err = profileFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = cpu.LoadConfig(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	processor, err = config.NewCpu()
	if err != nil {
		return fmt.Errorf("error processing config: %v", err)
	}
	statisticRequested := (*outputFileName != "")

	if *binaryFileName == "" {
		return fmt.Errorf("no program specified")
	}

	if statisticRequested {
		if *labelFileName == "" {
			return fmt.Errorf("a label file has to be specified")
		}

		if *percentageCutOff > 100 {
			return fmt.Errorf("%d is not a valid value for cutoff percentage", *percentageCutOff)
		}

		labels, err = acmeassembler.ParseLabelFile(*labelFileName)
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

		if err = profiler.DumpStatistics(processor.Mem, *outputFileName, labels, loadAddress, (loadAddress + progLen - 1), ctOff); err != nil {
			return fmt.Errorf("problem generating output file: %v", err)
		}
	}

	// picProc.Save("apfel.png")

	return nil
}
