package repository

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/platform/database"
)

type CommentRepository interface {
	Create(articleId int64, userId int64, body string) (int64, error)
	Delete(commentId int64) error
	Get(commentId int64) (*model.Comment, error)
	GetArticleComments(articleId int64) ([]*model.Comment, error)
}

func NewCommentRepository(db *database.DB) CommentRepository {
	return &CommentRepo{db}
}

type CommentRepo struct {
	db *database.DB
}

func (repo *CommentRepo) Create(articleId int64, userId int64, body string) (int64, error) {
	query := `INSERT INTO comment (article_id, author_id, body) VALUES(?, ?, ?)`
	result, err := repo.db.Exec(query, articleId, userId, body)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo *CommentRepo) Delete(commentId int64) error {
	query := `DELETE FROM comment WHERE id = ?`
	_, err := repo.db.Exec(query, commentId)
	return err
}

func (repo *CommentRepo) GetArticleComments(articleId int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	query := `SELECT * FROM comment WHERE article_id = ?`
	err := repo.db.Select(&comments, query, articleId)

	return comments, err
}

func (repo *CommentRepo) Get(commentId int64) (*model.Comment, error) {
	var comment model.Comment
	query := `SELECT * FROM comment WHERE id = ?`
	err := repo.db.Get(&comment, query, commentId)

	return &comment, err
}
