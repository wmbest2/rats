package project

import (
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"time"
)

type Project struct {
	Id        int64     `json:"-"`
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on,omitempty"`
}

func New(name string) (*Project, error) {
	project := &Project{Name: name}

	err := db.Conn.QueryRow(createProject, project.Name).Scan(&project.Id)
	if err != nil {
		return nil, err
	}

	api.GenerateToken(project)

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

// Be a TokenHolder
func (p *Project) Seed() string {
	return p.Name
}

func (p *Project) Type() api.TokenType {
	return api.ProjectToken
}

func (p *Project) Identifier() int64 {
	return p.Id
}
