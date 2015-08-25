package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"fmt"
	"github.com/wmbest2/rats/db"
	"time"
)

type TokenType int

const (
	UserToken TokenType = iota
	ProjectToken
)

type Token struct {
	Id             int64     `json:"id"`
	Token          string    `json:"-"`
	TokenEncrypted string    `json:"token"`
	Type           TokenType `json:"type"`
	ParentId       int64     `json:"parent_id"`
	CreatedOn      time.Time `json:"created_on,omitempty"`
}

type TokenHolder interface {
	Seed() string
	Type() TokenType
	Identifier() int64
}

var (
	DefaultCost = bcrypt.DefaultCost
)

func GenerateToken(holder TokenHolder) (string, error) {
	seed := fmt.Sprintf("%s%s%v", holder.Seed(), time.Now().UnixNano())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(seed)))

	token, err := bcrypt.GenerateFromPassword([]byte(hash), DefaultCost)
	if err != nil {
		return "", err
	}

	oldToken, _ := FetchToken(holder)
	if oldToken == "" {
		_, err = db.Conn.Exec(createToken, holder.Type(), holder.Identifier(), hash, token)
	} else {
		_, err = db.Conn.Exec(updateToken, holder.Identifier(), hash, token)
	}
	return hash, err
}

func FindToken(holder TokenHolder) (*Token, error) {
	var token Token
	println(holder.Identifier())
	err := db.Conn.QueryRow(findToken, holder.Identifier()).Scan(&token.Id,
		&token.Token,
		&token.TokenEncrypted,
		&token.Type,
		&token.ParentId,
		&token.CreatedOn)

	if err != nil {
		return nil, err
	}
	return &token, nil
}

func FetchToken(holder TokenHolder) (string, error) {
	var token Token
	err := db.Conn.QueryRow(findToken, holder.Identifier()).Scan(&token.Id,
		&token.Token,
		&token.TokenEncrypted,
		&token.Type,
		&token.ParentId,
		&token.CreatedOn)

	if err != nil {
		return "", err
	}
	return token.Token, nil
}
