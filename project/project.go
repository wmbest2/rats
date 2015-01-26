package project

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"fmt"
	"github.com/wmbest2/rats/db"
	"time"
)

var (
	DefaultCost = bcrypt.DefaultCost
)

type Project struct {
	Id        int64     `json:"-"`
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on,omitempty"`
}

type ProjectToken struct {
	Token          string    `json:"-"`
	TokenEncrypted string    `json:"-"`
	Project        int64     `json:"-"`
	CreatedOn      time.Time `json:"created_on,omitempty"`
}

func New(name string) (*Project, error) {
	project := &Project{Name: name}

	err := db.Conn.QueryRow(createProject, project.Name).Scan(&project.Id)
	if err != nil {
		return nil, err
	}

	project.GenerateToken()

	return project, err
}

func Find(name string) (*Project, error) {
	project := Project{}
	err := db.Conn.QueryRow(findProject, name).Scan(&project.Id, &project.Name, &project.CreatedOn)

	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *Project) FetchToken() (string, error) {
	var token ProjectToken
	err := db.Conn.QueryRow(findToken, p.Id).Scan(&token.Token, &token.TokenEncrypted, &token.Project, &token.CreatedOn)

	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (p *Project) GenerateToken() (string, error) {
	seed := fmt.Sprintf("%s%s%v", p.Name, time.Now().UnixNano())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(seed)))

	token, err := bcrypt.GenerateFromPassword([]byte(hash), DefaultCost)
	if err != nil {
		return "", err
	}

	oldToken, _ := p.FetchToken()
	if oldToken == "" {
		_, err = db.Conn.Exec(createProjectToken, p.Id, hash, token)
	} else {
		_, err = db.Conn.Exec(updateProjectToken, p.Id, hash, token)
	}
	return hash, err
}
