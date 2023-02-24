package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
)

type IterProcFunc func(testCaseName string) error

type CaseRepo interface {
	IterateTestCases(iterProcessor IterProcFunc) (uint, error)
	Get(caseName string) (*TestCase, error)
	New(caseName string, t *TestCase, createDriver bool) error
	GetScriptPath() string
}

func NewCaseRepo(testDir string) (CaseRepo, error) {
	return &simpleCaseRepo{
		testDir: testDir,
	}, nil
}

type simpleCaseRepo struct {
	testDir string
}

func (s *simpleCaseRepo) GetScriptPath() string {
	return s.testDir
}

func (s *simpleCaseRepo) New(caseName string, t *TestCase, createDriver bool) error {
	scriptPath := path.Join(s.testDir, t.TestScript)
	testDriverPath := path.Join(s.testDir, t.TestDriverSource)
	jsonPath := path.Join(s.testDir, caseName+".json")

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
	defer func() { f.Close() }()

	if createDriver {
		f2, err := os.Create(testDriverPath)
		if err != nil {
			return fmt.Errorf("unable to create test driver '%s'", testDriverPath)
		}
		defer func() { f2.Close() }()
	}

	return nil
}

func (s *simpleCaseRepo) Get(caseName string) (*TestCase, error) {
	caseFileName := path.Join(s.testDir, caseName)
	return NewTestCaseFromFile(caseFileName)
}

func (s *simpleCaseRepo) IterateTestCases(iterProcessor IterProcFunc) (uint, error) {
	r := regexp.MustCompile(`^(.+)\.json$`)

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
			err := iterProcessor(j)
			if err != nil {
				return testCount, err
			}

			testCount++
		}
	}

	return testCount, nil
}
