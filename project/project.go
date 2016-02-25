package project

import (
	"fmt"
	"github.com/wmbest2/rats/db"
	"github.com/wmbest2/rats/namedaccess"
	"time"
)

type Project struct {
	Id        int64     `json:"-"`
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on,omitempty"`
}

func New(name string, createToken bool) (*Project, error) {
	project := &Project{Name: name}

	err := db.Conn.QueryRow(createProject, project.Name).Scan(&project.Id)
	if err != nil {
		return nil, err
	}

	if createToken {
		namedaccess.NewWithProject(fmt.Sprintf("%s project token", name), project.Id)
	}

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
