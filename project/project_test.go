package project

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"github.com/wmbest2/rats/namedaccess"
	"testing"
)

func TestProjectCreation(t *testing.T) {

	Convey("Given a fresh database", t, func() {

		Convey("A Project token table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM runs`)
			So(err, ShouldBeNil)

			_, err = db.Conn.Exec(`DELETE FROM api_tokens`)
			So(err, ShouldBeNil)
		})

		Convey("A Project table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM projects`)
			So(err, ShouldBeNil)
		})

		Convey("When a project is created", func() {
			project, err := New("test_project", true)
			So(project, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Then a token should also be generated", func() {
				namedAccess, err := namedaccess.FindByProject(project.Id)
				So(err, ShouldBeNil)

				_, err = api.FetchToken(namedAccess)
				So(err, ShouldBeNil)
			})

		})

		Convey("When a project with a duplicate name is created", func() {
			project, err := New("test_project", false)
			So(project, ShouldBeNil)
			So(err, ShouldNotBeNil)

			Convey("Then the database should only contain the first project", func() {
				rows, err := db.Conn.Query(findProject, "test_project")
				So(err, ShouldBeNil)

				row_count := 0
				for rows.Next() {
					row_count++
				}

				So(1, ShouldEqual, row_count)
			})

		})
	})

	Convey("Given a project not yet created", t, func() {
		Convey("When finding the missing project", func() {
			Convey("Then an error should occur", func() {
				project, err := Find("tn_checker")
				So(err, ShouldNotBeNil)
				So(project, ShouldBeNil)
			})
		})
	})

}
