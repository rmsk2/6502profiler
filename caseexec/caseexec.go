package caseexec

import (
	"6502profiler/cpu"
	"6502profiler/emuconfig"
	"6502profiler/memory"
	"6502profiler/verifier"
	"fmt"
	"io"
	"os"
	"strings"
)

type CaseExec struct {
	cpuProv            emuconfig.CpuProvider
	originalProv       emuconfig.CpuProvider
	asmProv            emuconfig.AsmProvider
	repo               verifier.CaseRepo
	verboseFlag        bool
	Outf               io.Writer
	ReportAsmError     AsmErrorReporter
	ReportSummary      SummaryReporter
	SubCaseReporter    verifier.SubcaseProcessor
	ReportTestInfo     TestInfoReporter
	CurrentCpu         *cpu.CPU6502
	trapAddress        uint16
	placeholderWrapper *memory.PlaceholderWrapper
}

type AsmErrorReporter func(errMsg string)
type SummaryReporter func()
type TestInfoReporter func(string, *verifier.TestCase)

func NewCaseExec(c emuconfig.CpuProvider, a emuconfig.AsmProvider, repo verifier.CaseRepo, v bool) *CaseExec {
	res := CaseExec{
		cpuProv:            nil,
		originalProv:       c,
		asmProv:            a,
		repo:               repo,
		verboseFlag:        v,
		Outf:               os.Stdout,
		placeholderWrapper: nil,
		trapAddress:        0,
	}

	res.cpuProv = NewWrapperCpuProvider(c, &res)
	res.ReportAsmError = res.printAsmError
	res.ReportSummary = res.printSummary
	res.SubCaseReporter = res.printSubcaseInfo
	res.ReportTestInfo = res.printTestInfo

	return &res
}

func (t *CaseExec) SetTrapAddress(a uint16) {
	t.trapAddress = a
}

func (t *CaseExec) LoadAndExecuteCase(testCaseName string) error {
	caseFileName := testCaseName

	if !strings.HasSuffix(caseFileName, verifier.TestCaseExtension) {
		caseFileName += verifier.TestCaseExtension
	}

	testCase, err := t.repo.Get(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	return t.ExecuteCase(testCaseName, testCase)
}

func (t *CaseExec) ExecuteCase(testCaseName string, testCase *verifier.TestCase) error {
	t.placeholderWrapper = nil

	// This sets t.placeholderWrapper if a trap address is desired
	cpu, err := t.cpuProv.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test case: %v", err)
	}
	defer func() { cpu.Mem.Close() }()

	t.CurrentCpu = cpu

	assembler := t.asmProv.GetAssembler()
	var subcaseProc verifier.SubcaseProcessor = nil

	if t.verboseFlag {
		subcaseProc = t.SubCaseReporter
	}

	t.ReportTestInfo(testCaseName, testCase)

	err = testCase.Execute(cpu, assembler, t.repo.GetScriptPath(), subcaseProc, t.placeholderWrapper)
	if err != nil {
		errMsg := assembler.GetErrorMessage()
		if errMsg != "" {
			t.ReportAsmError(errMsg)
		}
		return fmt.Errorf("test case '%s' failed: %v", testCase.Name, err)
	}

	t.ReportSummary()

	return nil
}

func (t *CaseExec) printTestInfo(testCaseName string, testCase *verifier.TestCase) {
	if t.verboseFlag {
		fmt.Fprintln(t.Outf, "--------------------------------------------")
		fmt.Fprintf(t.Outf, "Executing test case '%s'\n", testCase.Name)
		fmt.Fprintf(t.Outf, "Test case file: %s\n", testCaseName)
		fmt.Fprintf(t.Outf, "Test script: %s\n", testCase.TestScript)
		fmt.Fprintf(t.Outf, "Test driver: %s\n", testCase.TestDriverSource)
	} else {
		fmt.Fprintf(t.Outf, "Executing test case '%s' ... ", testCase.Name)
	}
}

func (t *CaseExec) printSubcaseInfo(i uint, numIter uint) {
	fmt.Fprintf(t.Outf, "Executing subcase %d of %d (%d clock cycles already used)\n", i+1, numIter, t.CurrentCpu.NumCycles())
}

func (t *CaseExec) printSummary() {
	if t.verboseFlag {
		fmt.Fprintf(t.Outf, "Clock cycles used: %d\n", t.CurrentCpu.NumCycles())
		fmt.Fprintln(t.Outf, "Test result: OK")
	} else {
		fmt.Fprintf(t.Outf, "(%d clock cycles) OK\n", t.CurrentCpu.NumCycles())
	}
}

func (t *CaseExec) printAsmError(errMsg string) {
	fmt.Fprintln(t.Outf, errMsg)
}

func (t *CaseExec) ExecuteSetupProgram(setupPrgName string) error {
	asm := t.asmProv.GetAssembler()

	binaryName, err := asm.Assemble(setupPrgName)
	if err != nil {
		errMsg := asm.GetErrorMessage()
		if errMsg != "" {
			t.ReportAsmError(errMsg)
		}

		return fmt.Errorf("unable to setup tests: %v", err)
	}

	cpu, err := t.originalProv.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test setup: %v", err)
	}

	_, _, err = cpu.LoadAndRun(binaryName)
	if err != nil {
		return fmt.Errorf("unable to create cpu for test setup: %v", err)
	}

	res, err := newSnapshotProvider(cpu, t)
	if err != nil {
		return fmt.Errorf("unable to perform global test setup: %v", err)
	}

	t.cpuProv = res

	return nil
}

type snapshotCpuProvider struct {
	cpu *cpu.CPU6502
	ce  *CaseExec
	p   *memory.PlaceholderWrapper
}

func newSnapshotProvider(cpu *cpu.CPU6502, c *CaseExec) (emuconfig.CpuProvider, error) {
	cpu.Reset()
	cpu.Mem.TakeSnapshot()
	var placeholder *memory.PlaceholderWrapper = nil

	if c.trapAddress != 0 {
		placeholder := memory.NewPlaceholderWrapper(cpu.Mem, c.trapAddress)
		cpu.Mem = placeholder.Wrapper
	}

	return &snapshotCpuProvider{
		cpu: cpu,
		ce:  c,
		p:   placeholder,
	}, nil
}

func (c *snapshotCpuProvider) NewCpu() (*cpu.CPU6502, error) {
	c.cpu.Mem.RestoreSnapshot()
	c.cpu.Reset()
	c.ce.placeholderWrapper = c.p

	return c.cpu, nil
}

type wrapperCpuProvider struct {
	originalProv emuconfig.CpuProvider
	caseExec     *CaseExec
}

func NewWrapperCpuProvider(o emuconfig.CpuProvider, c *CaseExec) *wrapperCpuProvider {
	return &wrapperCpuProvider{
		originalProv: o,
		caseExec:     c,
	}
}

func (w *wrapperCpuProvider) NewCpu() (*cpu.CPU6502, error) {
	cpu, err := w.originalProv.NewCpu()
	if err != nil {
		return nil, err
	}

	if w.caseExec.trapAddress != 0 {
		w.caseExec.placeholderWrapper = memory.NewPlaceholderWrapper(cpu.Mem, w.caseExec.trapAddress)
		cpu.Mem = w.caseExec.placeholderWrapper.Wrapper
	}

	return cpu, nil
}
