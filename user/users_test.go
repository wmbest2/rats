package user

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"testing"
)

func TestProjectCreation(t *testing.T) {

	Convey("Given a fresh database", t, func() {

		Convey("A User table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM users WHERE username <> 'admin'`)
			So(err, ShouldBeNil)
		})

		Convey("A Default Admin user should be created", func() {
			_, err := db.Conn.Query(findUser, "admin")
			So(err, ShouldBeNil)
		})

		Convey("When a user is created", func() {
			user, err := New("new_user", "password", false)
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)

			Convey("Then a token should also be generated", func() {
				_, err := api.FetchToken(user)
				So(err, ShouldBeNil)
			})

		})

		Convey("When a user with a duplicate name is created", func() {
			project, err := New("new_user", "password", false)
			So(project, ShouldBeNil)
			So(err, ShouldNotBeNil)

			Convey("Then the database should only contain the original user", func() {
				rows, err := db.Conn.Query(findUser, "new_user")
				So(err, ShouldBeNil)

				row_count := 0
				for rows.Next() {
					row_count++
				}

				So(1, ShouldEqual, row_count)
			})

		})
	})
}
