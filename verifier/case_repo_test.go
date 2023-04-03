package verifier

import (
	"fmt"
	"testing"
)

type testRepo struct {
	testCases map[string]*TestCase
}

func NewTestRepo() *testRepo {
	return &testRepo{
		testCases: map[string]*TestCase{},
	}
}

func (t *testRepo) IterateTestCases(iterProcessor IterProcFunc) (uint, error) {
	var testCount uint = 0

	for i, j := range t.testCases {
		err := iterProcessor(i, j)
		if err != nil {
			return 0, err
		}
	}
	return testCount, nil
}

func (t *testRepo) Get(caseName string) (*TestCase, error) {
	if tc, ok := t.testCases[caseName]; ok {
		return tc, nil
	} else {
		return nil, fmt.Errorf("error loading test case")
	}
}

func (t *testRepo) Del(caseName string) error {
	delete(t.testCases, caseName)

	return nil
}

func (t *testRepo) Add(caseName string, tc *TestCase, createDriver bool) error {
	t.testCases[caseName] = tc

	return nil
}

func (t *testRepo) GetScriptPath() string {
	return ""
}

func TestStat(t *testing.T) {
	repo := NewTestRepo()
	c1 := NewTestCase("Test case 1", "test1")
	c2 := NewTestCaseWithDriver("Test case 2", "test2", "test1.a")
	c3 := NewTestCase("Test case 3", "test3")

	repo.Add("test1", c1, true)
	repo.Add("test2", c2, false)
	repo.Add("test3", c3, true)

	tc, asmU, luaU, err := statCase(repo, "test1")
	if err != nil {
		t.Fatal(err)
	}

	if asmU == true {
		t.Fatal("asm in case 1 is not unique")
	}

	if luaU == false {
		t.Fatal("Lua in case 1 is unique")
	}

	if tc.TestDriverSource != "test1.a" {
		t.Fatal("Wrong  driver name in case 1")
	}

	tc, asmU, luaU, err = statCase(repo, "test3")
	if err != nil {
		t.Fatal(err)
	}

	if asmU == false {
		t.Fatal("asm in case 3 is unique")
	}

	if luaU == false {
		t.Fatal("Lua in case 3 is unique")
	}

	if tc.TestDriverSource != "test3.a" {
		t.Fatal("Wrong  driver name in case 3")
	}

	tc, asmU, luaU, err = statCase(repo, "test2")
	if err != nil {
		t.Fatal(err)
	}

	if asmU == true {
		t.Fatal("asm in case 2 is not unique")
	}

	if luaU == false {
		t.Fatal("Lua in case 2 is unique")
	}

	if tc.TestDriverSource != "test1.a" {
		t.Fatal("Wrong  driver name in case 2")
	}
}

func TestStatLuaNotUnique(t *testing.T) {
	repo := NewTestRepo()
	c1 := TestCase{
		Name:             "Test case 1",
		TestDriverSource: "test1.a",
		TestScript:       "global.lua",
	}

	c2 := TestCase{
		Name:             "Test case 2",
		TestDriverSource: "test2.a",
		TestScript:       "global.lua",
	}

	repo.Add("test1", &c1, true)
	repo.Add("test2", &c2, true)

	tc, asmU, luaU, err := statCase(repo, "test1")
	if err != nil {
		t.Fatal(err)
	}

	if asmU == false {
		t.Fatal("asm in case 1 is unique")
	}

	if luaU == true {
		t.Fatal("Lua in case 1 is not unique")
	}

	if tc.TestDriverSource != "test1.a" {
		t.Fatalf("Wrong  driver name in case 1: %s", tc.TestDriverSource)
	}

	tc, asmU, luaU, err = statCase(repo, "test2")
	if err != nil {
		t.Fatal(err)
	}

	if asmU == false {
		t.Fatal("asm in case 2 is unique")
	}

	if luaU == true {
		t.Fatal("Lua in case 2 is not unique")
	}

	if tc.TestDriverSource != "test2.a" {
		t.Fatal("Wrong  driver name in case 2")
	}
}
