package runs

import (
	"time"

	"github.com/wmbest2/rats/db"
	"github.com/wmbest2/rats/test"
)

func NewRun(project int64, token int64) (*test.TestRun, error) {
	run := &test.TestRun{TokenId: token, ProjectId: project, Timestamp: time.Now()}

	err := db.Conn.QueryRow(createRun, token, project, run.Timestamp).Scan(&run.Id)
	if err != nil {
		return nil, err
	}
	return run, err
}

func SaveRun(r *test.TestRun) error {
	_, err := db.Conn.Exec(saveRun, r.Id, r.CommitId, r.Message, r.Time, r.Success)

	// recurse to store all suites/cases etc.
	return err
}

func FindTestRun(runId int64, loadAll bool) (*test.TestRun, error) {
	run := &test.TestRun{}
	err := db.Conn.QueryRow(findRun, runId).Scan(&run.Id,
		&run.TokenId,
		&run.ProjectId,
		&run.CommitId,
		&run.Message,
		&run.Timestamp,
		&run.Time,
		&run.Success)
	if err != nil {
		return nil, err
	}

	if loadAll {
		run.TestSuites, err = FindTestSuites(run, loadAll)
	}

	return run, err
}

func FindArtifacts(r *test.TestRun, recursive bool) ([]test.Artifact, error) {
	return nil, nil
}

func FindTestSuites(r *test.TestRun, recursive bool) ([]test.TestSuite, error) {
	return nil, nil
}

//func (s *TestSuite) FindTestCases(recusive bool) ([]TestCase, error) {

//return nil, nil
//}

//func (s *TestSuite) FindProperties() ([]Property, error) {

//return nil, nil
//}

//func (s *TestSuite) FindStacks(t StackType) ([]string, error) {

//return nil, nil
//}

//func (s *TestSuite) SaveSuite() error {

//return nil
//}

//func (s *TestSuite) AddCase(test *TestCase) error {

//return nil
//}
