package repository

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/platform/database"
)

type ArticleRepository interface {
	Create(article *model.Article) (int64, error)
	Get(ID int64, userId *int64) (*model.Article, error)
	GetBySlug(slug string, userId *int64) (*model.Article, error)
	GetArticleTags(articleId int64) ([]string, error)
	CreateArticleTag(articleId int64, tag string) (int64, error)
	DeleteArticleTag(articleId int64, tag string) error
	Update(ID int64, article *model.Article) error
	Delete(ID int64) error
	Favorited(articleId int64, userId int64) error
	Unfavorited(articleId int64, userId int64) error
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

func (repo ArticleRepo) Get(ID int64, userId *int64) (*model.Article, error) {
	article := model.Article{}
	query := `
		SELECT 
			article.*,
			(SELECT COUNT(id) FROM favorites_article WHERE favorites_article.article_id = article.id) AS favorites_count,
			(SELECT IF(COUNT(id) > 0, TRUE, FALSE) FROM favorites_article WHERE favorites_article.article_id = article.id and favorites_article.user_id = ?) AS favorited
		FROM article 
		WHERE article.id = ?;
	`
	err := repo.db.Get(&article, query, userId, ID)

	return &article, err
}

func (repo ArticleRepo) GetBySlug(slug string, userId *int64) (*model.Article, error) {
	article := model.Article{}
	query := `
		SELECT 
			article.*,
			(SELECT COUNT(id) FROM favorites_article WHERE favorites_article.article_id = article.id) AS favorites_count,
			(SELECT IF(COUNT(id) > 0, TRUE, FALSE) FROM favorites_article WHERE favorites_article.article_id = article.id and favorites_article.user_id = ?) AS favorited
		FROM article 
		WHERE article.slug = ?;
	`
	err := repo.db.Get(&article, query, userId, slug)

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

func (repo ArticleRepo) Favorited(articleId int64, userId int64) error {
	query := `INSERT INTO favorites_article (article_id, user_id) VALUES(?, ?)`
	_, err := repo.db.Exec(query, articleId, userId)
	return err
}

func (repo ArticleRepo) Unfavorited(articleId int64, userId int64) error {
	query := `DELETE FROM favorites_article WHERE article_id = ? and user_id = ?`
	_, err := repo.db.Exec(query, articleId, userId)
	return err
}
