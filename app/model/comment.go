package model

import "time"

type Comment struct {
	ID        int64      `db:"id" json:"id"`
	AuthorID  int64      `db:"author_id" json:"author_id"`
	ArticleID int64      `db:"article_id" json:"article_id"`
	Body      string     `db:"body" json:"body"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	Author    *Author    `json:"author"`
}

func NewComment() *Comment { return &Comment{} }

type CommentDTO struct {
	AuthorID int64  `json:"author_id" validate:"required,number"`
	Body     string `json:"body" validate:"required"`
}

type SaveCommentDTO struct {
	Comment CommentDTO `json:"comment`
}
