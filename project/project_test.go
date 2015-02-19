package project

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"testing"
)

func TestProjectCreation(t *testing.T) {

	Convey("Given a fresh database", t, func() {

		Convey("A Project token table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM api_tokens`)
			So(err, ShouldBeNil)
		})

		Convey("A Project table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM projects`)
			So(err, ShouldBeNil)
		})

		Convey("When a project is created", func() {
			project, err := New("test_project")
			So(project, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Then a token should also be generated", func() {
				_, err := api.FetchToken(project)
				So(err, ShouldBeNil)
			})

		})

		Convey("When a project with a duplicate name is created", func() {
			project, err := New("test_project")
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

func TestProjectCanOnlyHaveOneKey(t *testing.T) {

	Convey("Given an existing project token", t, func() {
		project, err := Find("check")
		if err != nil {
			project, err = New("check")
		}
		So(err, ShouldBeNil)
		So(project, ShouldNotBeNil)

		Convey("When a new token is generated", func() {
			project, err := Find("check")
			So(err, ShouldBeNil)
			So(project, ShouldNotBeNil)

			token, err := api.FetchToken(project)
			So(token, ShouldNotBeNil)
			So(err, ShouldBeNil)

			token_value, err := api.GenerateToken(project)
			So(token_value, ShouldNotEqual, token)

			Convey("Then it should replace the existing", func() {
				token, err := api.FetchToken(project)
				So(token_value, ShouldEqual, token)
				So(err, ShouldBeNil)
			})
		})
	})

}
