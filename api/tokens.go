package api

import (
	"crypto/md5"
	"fmt"
	"github.com/wmbest2/rats/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type TokenType int

const (
	UserToken TokenType = iota
	CIToken
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

func FindEncryptedToken(token string) (int64, error) {
	var id int64
	err := db.Conn.QueryRow(findEncryptedToken, token).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
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
