package repository

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/platform/database"
)

type ArticleRepository interface {
	Create(article *model.Article) (int64, error)
	Get(ID int64) (*model.Article, error)
	GetBySlug(slug string) (*model.Article, error)
	GetArticleTags(article_id int64) ([]string, error)
	CreateArticleTag(article_id int64, tag string) (int64, error)
	DeleteArticleTag(article_id int64, tag string) error
	Update(ID int64, article *model.Article) error
	Delete(ID int64) error
}

type ArticleRepo struct {
	db *database.DB
}

func NewArticleRepo(db *database.DB) ArticleRepository { return &ArticleRepo{db: db} }

func (repo ArticleRepo) Create(article *model.Article) (int64, error) {
	query := `INSERT INTO article (author_id, title, slug, description, body) VALUES(?, ?, ?, ?, ?)`
	result, err := repo.db.Exec(query, article.AuthorID, article.Title, article.Slug, article.Description, article.Body)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo ArticleRepo) Get(ID int64) (*model.Article, error) {
	article := model.Article{}
	query := `SELECT * FROM article WHERE id = ?`
	err := repo.db.Get(&article, query, ID)

	return &article, err
}

func (repo ArticleRepo) GetBySlug(slug string) (*model.Article, error) {
	article := model.Article{}
	query := `SELECT * FROM article WHERE slug = ?`
	err := repo.db.Get(&article, query, slug)

	return &article, err
}

func (repo ArticleRepo) GetArticleTags(article_id int64) ([]string, error) {
	var tags []string
	query := `SELECT tag_name FROM article_tags WHERE article_id = ?`
	err := repo.db.Select(&tags, query, article_id)

	return tags, err
}

func (repo ArticleRepo) Update(ID int64, article *model.Article) error {
	query := `UPDATE article SET title=?, description=?, body=?, updated_at=? WHERE id = ?`
	_, err := repo.db.Exec(query, article.Title, article.Description, article.Body, article.UpdatedAt, ID)
	return err
}

func (repo ArticleRepo) Delete(ID int64) error {
	query := `DELETE FROM article WHERE id = ?`
	_, err := repo.db.Exec(query, ID)
	return err
}

func (repo ArticleRepo) CreateArticleTag(article_id int64, tag string) (int64, error) {
	query := `INSERT INTO article_tags (article_id, tag_name) VALUES(?, ?)`
	result, err := repo.db.Exec(query, article_id, tag)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo ArticleRepo) DeleteArticleTag(article_id int64, tag string) error {
	query := `DELETE FROM article_tags WHERE article_id = ? and tag_name = ?`
	_, err := repo.db.Exec(query, article_id, tag)
	return err
}
