package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUser(user *model.User) (int, error) {
	var id int

	row := repo.db.QueryRow(`
		INSERT INTO tg_user(username, fullname, imageurl) 
		VALUES ($1, $2, $3) RETURNING id;`,
		user.Username, user.FullName, user.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (repo *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := repo.db.Get(&user, "SELECT * FROM tg_user WHERE username = $1;", username)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}

	return &user, nil
}

func (repo *UserRepo) GetUserByID(ID int) (*model.User, error) {
	var user model.User

	err := repo.db.Get(&user, "SELECT * FROM tg_user WHERE id = $1;", ID)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting user by ID: %w", err)
	}

	return &user, nil
}
