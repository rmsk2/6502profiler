package cpu

import (
	"fmt"
	"testing"
)

func TestAssemble(t *testing.T) {
	conf, err := LoadConfig("../config.json")
	if err != nil {
		t.Fatal(err)
	}

	p, err := conf.GetAcme().Assemble("fixed_test.a")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}
