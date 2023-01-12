package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type WebUserRepo struct {
	db *DB
}

func NewWebUserRepo(db *DB) *WebUserRepo {
	return &WebUserRepo{db: db}
}

func (repo *WebUserRepo) CreateWebUser(user *model.WebUser) error {
	_, err := repo.db.Exec(
		"INSERT INTO web_user(email, password) VALUES ($1, $2);",
		user.Email, user.Password,
	)
	if err != nil {
		return fmt.Errorf("create web user: %w", err)
	}

	return nil
}

func (repo *WebUserRepo) GetWebUserByID(id int) (*model.WebUser, error) {
	var user model.WebUser

	err := repo.db.Get(&user, "SELECT * FROM web_user WHERE id = $1;", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get web user by id: %w", err)
	}

	return &user, nil
}

func (repo *WebUserRepo) GetWebUserByEmail(email string) (*model.WebUser, error) {
	var user model.WebUser

	err := repo.db.Get(&user, "SELECT * FROM web_user WHERE email = $1;", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get web user by email: %w", err)
	}

	return &user, nil
}
