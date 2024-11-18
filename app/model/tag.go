package model

type Tag struct {
	Name string `db:"name" json:"name" validate:"required,lte=255,alphanum"`
}
