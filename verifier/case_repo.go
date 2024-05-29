package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

var emptyLuaScript string = `
function arrange()
	-- remember to call set_pc(load_address) if you use test iteration
end

function assert()
    return true, ""
end
`

type IterProcFunc func(testCaseName string, tCase *TestCase) error

type CaseRepo interface {
	IterateTestCases(iterProcessor IterProcFunc) (uint, error)
	Get(caseName string) (*TestCase, error)
	Add(caseName string, t *TestCase, createDriver bool) error
	Del(caseName string) error
	GetScriptPath() string
}

func NewCaseRepo(testDir string, defaultAsmDriver string) (CaseRepo, error) {
	return &simpleCaseRepo{
		testDir:       testDir,
		defaultDriver: defaultAsmDriver,
	}, nil
}

type simpleCaseRepo struct {
	testDir       string
	defaultDriver string
}

func (s *simpleCaseRepo) GetScriptPath() string {
	return s.testDir
}

func (s *simpleCaseRepo) Del(caseName string) error {
	if !strings.HasSuffix(caseName, TestCaseExtension) {
		caseName = caseName + TestCaseExtension
	}

	t, asmUnique, luaUnique, err := statCase(s, caseName)
	if err != nil {
		return fmt.Errorf("unable to delete case '%s': %v", caseName, err)
	}

	scriptPath := path.Join(s.testDir, t.TestScript)
	testDriverPath := path.Join(s.testDir, t.TestDriverSource)
	jsonPath := path.Join(s.testDir, caseName)

	err = os.Remove(jsonPath)
	if err != nil {
		return fmt.Errorf("unable to delete case '%s': %v", caseName, err)
	}

	if asmUnique {
		err = os.Remove(testDriverPath)
		if err != nil {
			return fmt.Errorf("unable to delete case '%s': %v", caseName, err)
		}
	}

	if luaUnique {
		err = os.Remove(scriptPath)
		if err != nil {
			return fmt.Errorf("unable to delete case '%s': %v", caseName, err)
		}
	}

	return nil
}

func (s *simpleCaseRepo) Add(caseName string, t *TestCase, createDriver bool) error {
	scriptPath := path.Join(s.testDir, t.TestScript)
	testDriverPath := path.Join(s.testDir, t.TestDriverSource)
	jsonPath := path.Join(s.testDir, caseName+TestCaseExtension)

	_, err := os.Stat(jsonPath)
	if err == nil {
		return fmt.Errorf("json file '%s' already exists", jsonPath)
	}

	_, err = os.Stat(scriptPath)
	if err == nil {
		return fmt.Errorf("script file '%s' already exists", scriptPath)
	}

	if createDriver {
		_, err = os.Stat(testDriverPath)
		if err == nil {
			return fmt.Errorf("test driver file '%s' already exists", testDriverPath)
		}
	}

	data, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to save testcase file %s: %v", jsonPath, err)
	}

	err = os.WriteFile(jsonPath, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to save config testcase file %s: %v", jsonPath, err)
	}

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to create lua script '%s'", scriptPath)
	}

	_, err = f.WriteString(emptyLuaScript)
	if err != nil {
		return fmt.Errorf("unable to initialize lua script '%s'", scriptPath)
	}

	defer func() { f.Close() }()

	if createDriver {
		f2, err := os.Create(testDriverPath)
		if err != nil {
			return fmt.Errorf("unable to create test driver '%s'", testDriverPath)
		}
		defer func() { f2.Close() }()

		_, err = f2.WriteString(s.defaultDriver)
		if err != nil {
			return fmt.Errorf("unable to initialize test driver '%s'", testDriverPath)
		}
	}

	return nil
}

func (s *simpleCaseRepo) Get(caseName string) (*TestCase, error) {
	caseFileName := path.Join(s.testDir, caseName)

	return NewTestCaseFromFile(caseFileName)
}

func statCase(s CaseRepo, caseName string) (tCase *TestCase, asmUnique bool, luaUnique bool, err error) {
	tCase, err = s.Get(caseName)
	if err != nil {
		return nil, false, false, err
	}

	asmCounter := map[string]int{}
	scriptCounter := map[string]int{}

	statIter := func(testCaseName string, caseData *TestCase) error {
		asmCounter[caseData.TestDriverSource] += 1
		scriptCounter[caseData.TestScript] += 1

		return nil
	}

	_, err = s.IterateTestCases(statIter)
	if err != nil {
		return nil, false, false, err
	}

	return tCase, asmCounter[tCase.TestDriverSource] == 1, scriptCounter[tCase.TestScript] == 1, nil
}

func (s *simpleCaseRepo) IterateTestCases(iterProcessor IterProcFunc) (uint, error) {
	r := regexp.MustCompile(fmt.Sprintf(`^(.+)\%s$`, TestCaseExtension))

	file, err := os.Open(s.testDir)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return 0, err
	}

	var testCount uint = 0

	for _, j := range names {
		if r.MatchString(j) {
			tCase, err := s.Get(j)
			if err != nil {
				return testCount, err
			}

			err = iterProcessor(j, tCase)
			if err != nil {
				return testCount, err
			}

			testCount++
		}
	}

	return testCount, nil
}
