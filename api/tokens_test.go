package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestTokenHolder struct {
	TType TokenType
}

func (t *TestTokenHolder) Seed() string {
	return "TEST TOKEN HOLDER"
}

func (t *TestTokenHolder) Type() TokenType {
	return t.TType
}

func (t *TestTokenHolder) Identifier() int64 {
	return 1
}

func TestDatabaseConnection(t *testing.T) {
	holder := &TestTokenHolder{}
	Convey("Given a token holder of type UserToken", t, func() {
		holder.TType = UserToken
		Convey("When I generate two tokens", func() {
			token1, err1 := GenerateToken(holder)
			token2, err2 := GenerateToken(holder)

			Convey("Then an Error should not occur", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})
			Convey("And they should be unique", func() {
				So(token1, ShouldNotEqual, token2)
			})
		})
	})
}
