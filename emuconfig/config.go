package emuconfig

import (
	"6502profiler/assembler"
	"6502profiler/cpu"
	"6502profiler/memory"
	"6502profiler/verifier"
	"encoding/json"
	"fmt"
	"os"
)

const UMul = 1
const SMul = 2
const UDiv = 4
const SDiv = 8

type Config struct {
	Model            string
	MemSpec          string
	IoMask           uint8
	IoAddrConfig     map[uint8]string
	PreLoad          map[uint16]string
	F256MCoprocFlags uint8
	F256MCoprocBase  uint16
	Ca65StartAddress uint16
	AsmType          string
	AcmeBinary       string
	AcmeSrcDir       string
	AcmeBinDir       string
	AcmeTestDir      string
}

type ConfParser func(cnf string) (memory.MemWrapper, bool)

var confParsers []ConfParser = []ConfParser{
	memory.NewFileProcFromConfig,
	memory.NewPicProcFromConfig,
	memory.NewStdOutProcessorFromConfig,
	memory.NewPrinterProcessorFromConfig,
}

const L16 = "Linear16K"
const L32 = "Linear32K"
const L48 = "Linear48K"
const L64 = "Linear64K"
const X16_512 = "XSixteen512K"
const X16_2048 = "XSixteen2048K"
const GEO_512 = "GeoRam_512K"
const GEO_2048 = "GeoRam_2048K"
const F256_512 = "F256_512K"
const F256_768 = "F256_768K"
const Proc6502 = "6502"
const Proc65C02 = "65C02"
const AsmDefault = ""
const AsmAcme = "acme"
const Asm64Tass = "64tass"
const AsmCa65 = "ca65"

const IllegalTrapAddress = 0
const Ca65DefaultLoadAddr = 0x0800

func NewConfigFromFile(fileName string) (*Config, error) {
	res := &Config{}

	allowedMemModels := map[string]bool{
		L16:      true,
		L32:      true,
		L48:      true,
		L64:      true,
		X16_512:  true,
		X16_2048: true,
		GEO_512:  true,
		GEO_2048: true,
		F256_512: true,
		F256_768: true,
	}

	allowedCpuModels := map[string]bool{
		Proc6502:  true,
		Proc65C02: true,
	}

	allowedAsmTypes := map[string]bool{
		AsmDefault: true,
		AsmAcme:    true,
		Asm64Tass:  true,
		AsmCa65:    true,
	}

	configData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to load config file %s: %v", fileName, err)
	}

	err = json.Unmarshal(configData, res)
	if err != nil {
		return nil, fmt.Errorf("unable to load config file %s: %v", fileName, err)
	}

	_, ok := allowedMemModels[res.MemSpec]
	if !ok {
		return nil, fmt.Errorf("unknown memory model: %v", res.MemSpec)
	}

	_, ok = allowedCpuModels[res.Model]
	if !ok {
		return nil, fmt.Errorf("unknown CPU model: %v", res.MemSpec)
	}

	_, ok = allowedAsmTypes[res.AsmType]
	if !ok {
		return nil, fmt.Errorf("unknown Assembler type: %v", res.AsmType)
	}

	return res, nil
}

func DefaultConfig() *Config {
	res := &Config{
		Model:            Proc6502,
		MemSpec:          L32,
		IoMask:           0,
		IoAddrConfig:     map[uint8]string{},
		PreLoad:          map[uint16]string{},
		F256MCoprocFlags: 0,
		F256MCoprocBase:  0xDE00,
		Ca65StartAddress: 0x0800,
		AsmType:          AsmAcme,
		AcmeBinary:       "acme",
		AcmeSrcDir:       "./",
		AcmeBinDir:       "./test/bin",
		AcmeTestDir:      "./test",
	}

	return res
}

func (c *Config) tryParseWrapperLine(line string) (memory.MemWrapper, bool) {
	var res memory.MemWrapper = nil
	ok := false

	for _, j := range confParsers {
		res, ok = j(line)
		if ok {
			return res, ok
		}
	}

	return res, ok
}

func (c *Config) Save(fileName string) error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to save config in file %s: %v", fileName, err)
	}

	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to save config in file %s: %v", fileName, err)
	}

	return nil
}

type CpuProvider interface {
	NewCpu() (*cpu.CPU6502, error)
}

type AsmProvider interface {
	GetAssembler() assembler.Assembler
}

type RepoProvider interface {
	GetCaseRepo() (verifier.CaseRepo, error)
}

func (c *Config) GetCaseRepo() (verifier.CaseRepo, error) {
	var repo verifier.CaseRepo

	repo, err := verifier.NewCaseRepo(c.AcmeTestDir)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (c *Config) GetAssembler() assembler.Assembler {
	switch {
	case c.AsmType == Asm64Tass:
		return assembler.NewTass64(c.AcmeBinary, c.AcmeSrcDir, c.AcmeBinDir, c.AcmeTestDir)
	case c.AsmType == AsmCa65:
		var loadAddress uint16 = Ca65DefaultLoadAddr
		if c.Ca65StartAddress != 0 {
			loadAddress = c.Ca65StartAddress
		}
		return assembler.NewCa65(c.AcmeBinary, c.AcmeSrcDir, c.AcmeBinDir, c.AcmeTestDir, loadAddress)
	default:
		return assembler.NewACME(c.AcmeBinary, c.AcmeSrcDir, c.AcmeBinDir, c.AcmeTestDir)
	}
}

func (c *Config) AddIoWrapper(mem memory.Memory) (memory.Memory, error) {
	var res memory.MemWrapper = nil
	ok := false
	baseAddress := uint16(c.IoMask) << 8
	var wrapper *memory.WrappingMemory = nil

	wrapper = memory.NewMemWrapper(mem, baseAddress)

	for i, j := range c.IoAddrConfig {
		res, ok = c.tryParseWrapperLine(j)
		if ok {
			res.SetBaseMem(mem)
			address := baseAddress | uint16(i)
			wrapper.AddWrapper(address, res)
		} else {
			return nil, fmt.Errorf("unable to process memory wrapper config: '%s'", j)
		}
	}

	mem = wrapper

	return mem, nil
}

func (c *Config) AddF256Func(mem memory.Memory) memory.Memory {
	var wrapper = memory.NewMemWrapper(mem, c.F256MCoprocBase)
	coproc := memory.NewUnsignedCoproc(mem, c.F256MCoprocBase)
	wrapperUsed := false

	if (c.F256MCoprocFlags & UMul) != 0 {
		coproc.RegisterUmul(wrapper)

		wrapperUsed = true
	}

	if (c.F256MCoprocFlags & UDiv) != 0 {
		coproc.RegisterUdiv(wrapper)

		wrapperUsed = true
	}

	if wrapperUsed {
		mem = wrapper
	}

	return mem
}

func (c *Config) PreloadRoms(cpu *cpu.CPU6502) error {
	for i, j := range c.PreLoad {
		data, err := os.ReadFile(j)
		if err != nil {
			return fmt.Errorf("unable to preload file: %v", err)
		}

		err = cpu.CopyToMem(data, i)
		if err != nil {
			return fmt.Errorf("unable to copy preloaded file '%s' to $%04X: %v", j, i, err)
		}
	}

	return nil
}

func (c *Config) NewCpu() (*cpu.CPU6502, error) {
	var model cpu.CpuModel = cpu.Model6502

	if c.Model != Proc6502 {
		model = cpu.Model65C02
	}

	cpu := cpu.New6502(model)
	var mem memory.Memory

	switch c.MemSpec {
	case L16:
		mem = memory.NewLinearMemory(16384)
	case L32:
		mem = memory.NewLinearMemory(32768)
	case L48:
		mem = memory.NewLinearMemory(49152)
	case X16_512:
		mem = memory.NewX16Memory(memory.X512K)
	case X16_2048:
		mem = memory.NewX16Memory(memory.X2048K)
	case GEO_512:
		mem = memory.NewNeoGeo(memory.NeoGeoRegisterPage+0xFE, 5)
	case GEO_2048:
		mem = memory.NewNeoGeo(memory.NeoGeoRegisterPage+0xFE, 7)
	case F256_512:
		mem = memory.NewF56JrMemory(false)
	case F256_768:
		mem = memory.NewF56JrMemory(true)
	default:
		mem = memory.NewLinearMemory(65536)
	}

	var err error

	mem = c.AddF256Func(mem)

	if len(c.IoAddrConfig) != 0 {
		mem, err = c.AddIoWrapper(mem)
		if err != nil {
			return nil, err
		}
	}

	cpu.Init(mem)

	err = c.PreloadRoms(cpu)
	if err != nil {
		return nil, err
	}

	return cpu, nil
}
