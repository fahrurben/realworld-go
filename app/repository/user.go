package repository

import (
	"database/sql"
	"fmt"
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/platform/database"
)

type UserRepository interface {
	Create(user *model.User) (int64, error)
	Get(ID int64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Exists(username, email string) (bool, error)
	Update(ID int64, user *model.User) error
	Delete(ID int64) error
	IsFollowing(user_id int64, follow_user_id int64) (bool, error)
	Follow(user_id int64, follow_user_id int64) error
	Unfollow(user_id int64, follow_user_id int64) error
}

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(user *model.User) (int64, error) {
	query := `INSERT INTO users (username, email, password) VALUES(?, ?, ?)`
	result, err := repo.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo *UserRepo) Get(ID int64) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM users WHERE id = ?`
	err := repo.db.Get(&user, query, ID)

	return &user, err
}

func (repo *UserRepo) GetByEmail(email string) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM users WHERE email = ?`
	err := repo.db.Get(&user, query, email)

	return &user, err
}

func (repo *UserRepo) GetByUsername(username string) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM users WHERE username = ?`
	err := repo.db.Get(&user, query, username)

	return &user, err
}

func (repo *UserRepo) Exists(username, email string) (bool, error) {
	findByUsername := len(username)
	findByEmail := len(email)

	query := `SELECT 1 FROM users`
	var placeholderValues []any
	if findByUsername > 0 && findByEmail > 0 {
		query = fmt.Sprintf("%s WHERE username = ? or email = ?", query)
		placeholderValues = append(placeholderValues, username, email)
	} else {
		if findByUsername > 0 {
			query = fmt.Sprint("%s WHERE username = ?", query)
			placeholderValues = append(placeholderValues, username)
		} else {
			query = fmt.Sprint("%s WHERE email = ?", query)
			placeholderValues = append(placeholderValues, email)
		}
	}

	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	var exists bool
	err := repo.db.QueryRow(query, placeholderValues...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, err
}

func (repo *UserRepo) Update(ID int64, user *model.User) error {
	query := `UPDATE users SET username=?, password=?, bio=?, image=? WHERE id = ?`
	_, err := repo.db.Exec(query, user.Username, user.Password, user.Bio, user.Bio, user.Image, ID)
	return err
}

func (repo *UserRepo) Delete(ID int64) error {
	query := `DELETE users WHERE id = ?`
	_, err := repo.db.Exec(query, ID)
	return err
}

func (repo *UserRepo) IsFollowing(user_id int64, follow_user_id int64) (bool, error) {
	var count int
	err := repo.db.Get(&count, "SELECT count(*) FROM following WHERE user_id = ? and follow_user_id = ?", user_id, follow_user_id)
	if count > 0 {
		return true, err
	}

	return false, err
}

func (repo *UserRepo) Follow(user_id int64, follow_user_id int64) error {
	query := "INSERT INTO following (user_id, follow_user_id) VALUES(?, ?)"
	_, err := repo.db.Exec(query, user_id, follow_user_id)
	return err
}

func (repo *UserRepo) Unfollow(user_id int64, follow_user_id int64) error {
	query := `DELETE FROM following WHERE user_id = ? and follow_user_id = ?`
	_, err := repo.db.Exec(query, user_id, follow_user_id)
	return err
}
