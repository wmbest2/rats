package runs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/project"
)

func TestCreateAndUpdateTestRun(t *testing.T) {

	project, _ := project.New("test_project_for_test_package")
	token, _ := api.FindToken(project)

	Convey("Given a newly created TestRun", t, func() {
		run, err := NewRun(project.Id, token.Id)
		Convey("When There are no errors", func() {
			So(err, ShouldBeNil)
			Convey("Then The Run should exist in the database", func() {
				fromDb, err := FindTestRun(run.Id, false)
				So(err, ShouldBeNil)
				So(fromDb.Id, ShouldEqual, run.Id)
			})

		})

		Convey("When Fields are updated and the Object is saved", func() {

			Convey("Then the data should be updated in the databse", func() {
			})

		})

	})

}
