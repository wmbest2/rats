package runs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/namedaccess"
	"github.com/wmbest2/rats/project"
)

func TestCreateAndUpdateTestRun(t *testing.T) {

	project, _ := project.New("test_project_for_test_package", true)
	access, _ := namedaccess.FindByProject(project.Id)
	token, _ := api.FindToken(access)

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

		Convey("When searching for a non-existent run ", func() {
			_, err := FindTestRun(-1, false)
			Convey("Then an Error should occure", func() {
				So(err, ShouldNotBeNil)
			})

		})

		Convey("When Fields are updated and the Object is saved", func() {

			Convey("Then the data should be updated in the database", func() {
			})

		})

	})

}
