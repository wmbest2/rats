package namedaccess

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"log"
	"testing"
)

func TestNamedAccessCreation(t *testing.T) {

	Convey("Given a fresh database", t, func() {

		Convey("A NamedAccess token table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM runs`)
			So(err, ShouldBeNil)

			_, err = db.Conn.Exec(`DELETE FROM api_tokens`)
			So(err, ShouldBeNil)
		})

		Convey("A NamedAccess table should be created", func() {
			_, err := db.Conn.Exec(`DELETE FROM namedaccess`)
			So(err, ShouldBeNil)
		})

		Convey("When a namedaccess is created", func() {
			namedaccess, err := New("test_namedaccess")
			So(namedaccess, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Then a token should also be generated", func() {
				_, err := api.FetchToken(namedaccess)
				So(err, ShouldBeNil)
			})

		})
	})

	Convey("Given a namedaccess not yet created", t, func() {
		Convey("When finding the missing namedaccess", func() {
			Convey("Then an error should occur", func() {
				namedaccess, err := Find("tn_checker")
				So(err, ShouldNotBeNil)
				So(namedaccess, ShouldBeNil)
			})
		})
	})

}

func TestNamedAccessCanOnlyHaveOneKey(t *testing.T) {

	Convey("Given an existing namedaccess token", t, func() {
		namedaccess, err := Find("check")
		if err != nil {
			namedaccess, err = New("check")
		}
		So(err, ShouldBeNil)
		So(namedaccess, ShouldNotBeNil)

		Convey("When a new token is generated", func() {
			namedaccess, err := Find("check")
			So(err, ShouldBeNil)
			So(namedaccess, ShouldNotBeNil)

			log.Printf("namedaccess id: %d\n", namedaccess.Id)

			token, err := api.FetchToken(namedaccess)
			So(token, ShouldNotEqual, "")
			So(err, ShouldBeNil)

			token_value, err := api.GenerateToken(namedaccess)
			So(token_value, ShouldNotEqual, token)

			Convey("Then it should replace the existing", func() {
				token, err := api.FetchToken(namedaccess)
				So(token_value, ShouldEqual, token)
				So(err, ShouldBeNil)
			})
		})
	})

}
