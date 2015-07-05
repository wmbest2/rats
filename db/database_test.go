package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {

	Convey("Given an initialized db", t, func() {

		Convey("When tables are selected", func() {

			Convey("Then an Error should not occur", func() {
				_, err := Conn.Exec(`SELECT count(*) FROM pg_tables`)
				So(err, ShouldBeNil)
			})

		})

		Convey("When calling ExecOrLog with an error", func() {
			Convey("Then the application shouldn't crash", func() {
				ExecOrLog("sdlkjfsl;")
			})
		})

	})
}
