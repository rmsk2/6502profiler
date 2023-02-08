package cpu

import (
	"6502profiler/memory"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Model        CpuModel
	MemSpec      string
	IoMask       uint8
	IoAddrConfig map[uint8]string
}

const L16 = "Linear16K"
const L32 = "Linear32K"
const L64 = "Linear64K"
const X16_512 = "XSixteen512K"
const X16_2048 = "XSixteen2048K"

func LoadConfig(fileName string) (*Config, error) {
	res := &Config{}

	allowedMemModels := map[string]bool{
		L16: true,
		L32: true,
		L64: true,
		//X16_512:  true,
		//X16_2048: true,
	}

	allowedCpuModels := map[CpuModel]bool{
		Model6502:  true,
		Model65C02: true,
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

	return res, nil
}

func EmptyConfig() *Config {
	res := &Config{
		Model:        Model6502,
		MemSpec:      L32,
		IoMask:       0,
		IoAddrConfig: map[uint8]string{},
	}

	return res
}

func (c *Config) SaveConfig(fileName string) error {
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

func (c *Config) NewCpu() (*CPU6502, error) {
	cpu := New6502(c.Model)
	var mem memory.Memory
	var wrapper *memory.WrappingMemory = nil

	switch c.MemSpec {
	case L16:
		mem = memory.NewLinearMemory(16384)
	case L32:
		mem = memory.NewLinearMemory(32768)
	default:
		mem = memory.NewLinearMemory(65536)
	}

	if len(c.IoAddrConfig) != 0 {
		wrapper = memory.NewMemWrapper(mem, uint16(c.IoMask)<<8)
		//for i, j := range c.IoAddrConfig {

		//}

		mem = wrapper
	}

	cpu.Init(mem)

	return cpu, nil
}