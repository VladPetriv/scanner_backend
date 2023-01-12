package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUser(user *model.User) (int, error) {
	row, err := repo.db.Exec(`
		INSERT INTO tg_user(username, fullname, image_url) VALUES ($1, $2, $3);`,
		user.Username, user.FullName, user.ImageURL,
	)
	if err != nil {
		return 0, fmt.Errorf("create tg user: %w", err)
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("create tg user: %w", err)
	}

	return int(id), nil
}

func (repo *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := repo.db.Get(&user, "SELECT * FROM tg_user WHERE username = $1;", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &user, nil
}

func (repo *UserRepo) GetUserByID(id int) (*model.User, error) {
	var user model.User

	err := repo.db.Get(&user, "SELECT * FROM tg_user WHERE id = $1;", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}
