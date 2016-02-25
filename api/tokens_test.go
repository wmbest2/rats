package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestTokenHolder struct {
	TType TokenType
	Id    int64
}

func (t *TestTokenHolder) Seed() string {
	return "TEST TOKEN HOLDER"
}

func (t *TestTokenHolder) Type() TokenType {
	return t.TType
}

func (t *TestTokenHolder) Identifier() int64 {
	return t.Id
}

func TestFailedEncryption(t *testing.T) {
	holder := &TestTokenHolder{Id: 1, TType: UserToken}
	Convey("Given an invalid default cost to bcrypt", t, func() {
		cost := DefaultCost
		DefaultCost = 10000
		Convey("When I generate a token", func() {
			token, err := GenerateToken(holder)

			Convey("Then an Error should occur", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the token should be nil", func() {
				So(token, ShouldEqual, "")
			})
		})
		DefaultCost = cost
	})
}

func TestDatabaseConnection(t *testing.T) {
	holder := &TestTokenHolder{Id: 1, TType: UserToken}
	Convey("Given a token holder of type UserToken", t, func() {
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

	Convey("Given a valid token", t, func() {
		token1, err1 := GenerateToken(holder)
		Convey("When I find the token", func() {
			token2, err2 := FindToken(holder)

			Convey("Then an Error should not occur", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})

			Convey("And the it should match the original", func() {
				So(token1, ShouldEqual, token2.Token)
			})
		})
	})

	Convey("Given a valid token", t, func() {
		holder.TType = UserToken
		token1, err1 := GenerateToken(holder)
		Convey("When I find the token", func() {
			token2, err2 := FindToken(holder)

			Convey("Then an Error should not occur", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
			})

			Convey("And the it should match the original", func() {
				So(token1, ShouldEqual, token2.Token)
			})
		})
		Convey("When I find the encrypted token", func() {
			token2, err2 := FindToken(holder)
			id, err3 := FindEncryptedToken(token2.TokenEncrypted)

			Convey("Then an Error should not occur", func() {
				So(err2, ShouldBeNil)
				So(err3, ShouldBeNil)
			})

			Convey("And the it should match the original", func() {
				So(id, ShouldEqual, token2.Id)
			})
		})
	})

	Convey("Given an invalid token", t, func() {
		temp := &TestTokenHolder{Id: 234}
		Convey("When I find the token", func() {
			token, err := FindToken(temp)

			Convey("Then an Error should be thrown", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And the token should be nil", func() {
				So(token, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid encrypted token", t, func() {
		Convey("When I find the token id", func() {
			id, err := FindEncryptedToken("faksdjfhlkashd")

			Convey("Then an Error should occur", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And the token id should be -1", func() {
				So(id, ShouldEqual, -1)
			})
		})
	})
}
