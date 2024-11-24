package repository

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/platform/database"
)

type TagRepository interface {
	Create(tag string) (int64, error)
	Get(name string) (*model.Tag, error)
	Delete(name string) error
}

type TagRepo struct {
	db *database.DB
}

func NewTagRepo(db *database.DB) TagRepository {
	return &TagRepo{db: db}
}

func (repo *TagRepo) Create(tag string) (int64, error) {
	query := `INSERT INTO tag (name) VALUES(?)`
	result, err := repo.db.Exec(query, tag)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo *TagRepo) Get(name string) (*model.Tag, error) {
	tag := model.Tag{}
	query := `SELECT * FROM tag WHERE name = ?`
	err := repo.db.Get(&tag, query, name)

	return &tag, err
}

func (repo *TagRepo) Delete(name string) error {
	query := `DELETE tag WHERE name = ?`
	_, err := repo.db.Exec(query, name)
	return err
}
