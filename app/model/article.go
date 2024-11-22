package model

import "time"

type Author struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

type Article struct {
	ID          int64      `db:"id" json:"id"`
	AuthorID    int64      `db:"author_id" json:"author_id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Description string     `db:"description" json:"description"`
	Body        string     `db:"body" json:"body"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	Author      *Author    `json:"author"`
	Tags        []string   `json:"tagList"`
}

func NewArticle() *Article { return &Article{} }

type ArticleDTO struct {
	AuthorID    string   `json:"author_id" validate:"required,number"`
	Title       string   `json:"title" validate:"required,lte=255"`
	Description string   `json:"description" validate:"required,lte=1000"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList"`
}

type SaveArticleDto struct {
	Article ArticleDTO `json:"article"`
}
