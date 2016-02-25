package namedaccess

import (
	"github.com/wmbest2/rats/api"
	"github.com/wmbest2/rats/db"
	"time"
)

type NamedAccess struct {
	Id        int64     `json:"-"`
	Name      string    `json:"name"`
	ProjectId int64     `json:"project_id"`
	CreatedOn time.Time `json:"created_on,omitempty"`
}

func insertNamedAccess(namedAccess *NamedAccess) error {
	return db.Conn.QueryRow(createNamedAccess, namedAccess.Name, namedAccess.ProjectId).Scan(&namedAccess.Id)
}

func New(name string) (*NamedAccess, error) {
	namedAccess := &NamedAccess{Name: name, ProjectId: -1}

	err := insertNamedAccess(namedAccess)
	if err != nil {
		return nil, err
	}

	api.GenerateToken(namedAccess)

	return namedAccess, err
}

func NewWithProject(name string, projectId int64) (*NamedAccess, error) {
	namedAccess := &NamedAccess{Name: name, ProjectId: projectId}

	err := insertNamedAccess(namedAccess)
	if err != nil {
		return nil, err
	}

	api.GenerateToken(namedAccess)

	return namedAccess, err
}

func Find(name string) (*NamedAccess, error) {
	namedAccess := NamedAccess{}
	err := db.Conn.QueryRow(findNamedAccess, name).Scan(&namedAccess.Id, &namedAccess.Name, &namedAccess.ProjectId, &namedAccess.CreatedOn)

	if err != nil {
		return nil, err
	}
	return &namedAccess, nil
}

func FindByProject(project int64) (*NamedAccess, error) {
	namedAccess := NamedAccess{}
	err := db.Conn.QueryRow(findNamedAccessByProject, project).Scan(&namedAccess.Id, &namedAccess.Name, &namedAccess.ProjectId, &namedAccess.CreatedOn)

	if err != nil {
		return nil, err
	}
	return &namedAccess, nil
}

// Be a TokenHolder
func (p *NamedAccess) Seed() string {
	return p.Name
}

func (p *NamedAccess) Type() api.TokenType {
	return api.CIToken
}

func (p *NamedAccess) Identifier() int64 {
	return p.Id
}
