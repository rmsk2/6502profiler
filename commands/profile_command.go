package commands

import (
	"6502profiler/cpu"
	"6502profiler/emuconfig"
	"6502profiler/luabridge"
	"6502profiler/memory"
	"6502profiler/profiler"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	lua "github.com/yuin/gopher-lua"
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

func LoadAndRunBinary(processor *cpu.CPU6502, binaryFileName *string, trapAddress *uint, trapScript *string, silent bool) (uint16, uint16, error) {
	loadAddress, progLen, err := processor.Load(*binaryFileName)
	if err != nil {
		return 0, 0, fmt.Errorf("%v", err)
	}

	var L *lua.LState

	if !silent {
		fmt.Printf("Program loaded to address $%04x\n", loadAddress)
	}

	if *trapAddress != emuconfig.IllegalTrapAddress {
		if *trapScript == "" {
			return 0, 0, fmt.Errorf("no Lua script specified")
		}

		L = lua.NewState()
		defer L.Close()

		trapAddr := (uint16)(*trapAddress)

		if !silent {
			fmt.Printf("Using trap address $%x\n", trapAddr)
		}

		baseMem := processor.Mem

		wrapperMem := memory.NewMemWrapper(baseMem, 0xFF00&trapAddr)
		trapProc, err := luabridge.NewTrapProcessor(L, *trapScript, processor, loadAddress, progLen, *binaryFileName+".ident")
		if err != nil {
			return 0, 0, fmt.Errorf("%v", err)
		}

		wrapperMem.AddSpecialWriteAddress(trapAddr, trapProc.Write)
		processor.Mem = wrapperMem
		defer func() {
			_ = trapProc.Ctx.CallCleanup()
			// Remove memory wrapper, because the trap adddress will not work after the Lua
			// state has been Closed.
			processor.Mem = baseMem
		}()
	}

	err = processor.Run(loadAddress)
	if err != nil {
		return 0, 0, fmt.Errorf("a problem occurred: %v", err)
	}

	return loadAddress, progLen, nil
}

func RunCommand(arguments []string) error {
	var config *emuconfig.Config = emuconfig.DefaultConfig()
	var processor *cpu.CPU6502
	var err error = nil

	runFlags := flag.NewFlagSet("6502profiler run", flag.ContinueOnError)
	binaryFileName := runFlags.String("prg", "", "Path to the program to run")
	configName := runFlags.String("c", "", "Config file name")
	dumpFlag := runFlags.String("dump", "", "Dump memory after program has stopped. Format 'startaddr:len'")
	trapAddress := runFlags.Uint("trapaddr", emuconfig.IllegalTrapAddress, "Address to use for triggering a trap")
	trapScript := runFlags.String("lua", "", "Lua script to call when trap is triggered")
	silent := runFlags.Bool("silent", false, "Do not print additional info")

	if err = runFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = emuconfig.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	processor, err = config.NewCpu()
	if err != nil {
		return fmt.Errorf("error processing config: %v", err)
	}

	if *binaryFileName == "" {
		return fmt.Errorf("no program specified")
	}

	if _, _, err := parseDumpParams(*dumpFlag); err != nil {
		return err
	}

	_, _, err = LoadAndRunBinary(processor, binaryFileName, trapAddress, trapScript, *silent)
	if err != nil {
		return err
	}

	if !*silent {
		fmt.Printf("Program ran for %d clock cycles\n", processor.NumCycles())
	}

	err = DumpMemory(*dumpFlag, processor)
	if err != nil {
		return err
	}

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
	trapAddress := profileFlags.Uint("trapaddr", emuconfig.IllegalTrapAddress, "Address to use for triggering a trap")
	trapScript := profileFlags.String("lua", "", "Lua script to call when trap is triggered")
	silent := profileFlags.Bool("silent", false, "Do not print additional info")

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

	statisticRequested := (*outputFileName != "")

	if *binaryFileName == "" {
		return fmt.Errorf("no program specified")
	}

	if _, _, err := parseDumpParams(*dumpFlag); err != nil {
		return err
	}

	if statisticRequested {
		labels = map[uint16][]string{}

		if *labelFileName != "" {
			labels, err = assembler.ParseLabelFile(*labelFileName)
			if err != nil {
				return fmt.Errorf("a problem occurred: %v", err)
			}
		}

		if *percentageCutOff > 100 {
			return fmt.Errorf("%d is not a valid value for cutoff percentage", *percentageCutOff)
		}

		p = float64(*percentageCutOff) / 100.0
	}

	loadAddress, progLen, err := LoadAndRunBinary(processor, binaryFileName, trapAddress, trapScript, *silent)
	if err != nil {
		return err
	}

	if !*silent {
		fmt.Printf("Program ran for %d clock cycles\n", processor.NumCycles())
	}

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
