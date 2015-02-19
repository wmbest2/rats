package user

import (
	"crypto/md5"
	"fmt"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
)

var (
	DefaultCost = bcrypt.DefaultCost
)

type User struct {
	Id        int64     `json:"-"`
	Username  string    `json:"username"`
	Password  *string   `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedOn time.Time `json:"created_on,omitempty"`
}

// Password may not be nil unless oauth?
func New(name string, password string, isAdmin bool) (*User, error) {
	passwordhash := EncryptPassword(password)
	user := &User{Username: name, Password: &passwordhash, IsAdmin: isAdmin}

	err := db.Conn.QueryRow(createUser, user.Username, user.Password, user.IsAdmin).Scan(&user.Id)
	if err != nil {
		return nil, err
	}

	api.GenerateToken(user)

	return user, err
}

func EncryptPassword(password string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	token, _ := bcrypt.GenerateFromPassword([]byte(hash), DefaultCost)
	return string(token)
}

func FindUserByOauth(service string, id string) (*User, error) {
	return nil, nil
}

func (u *User) AddOauth(service string, id string) error {
	return nil
}

// Be a TokenHolder
func (u *User) Seed() string {
	return u.Username
}

func (u *User) Type() api.TokenType {
	return api.UserToken
}

func (u *User) Identifier() int64 {
	return u.Id
}
