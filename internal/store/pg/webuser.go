package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

var ErrWebUserNotFound = errors.New("web user not found")

type WebUserRepo struct {
	db *DB
}

func NewWebUserRepo(db *DB) *WebUserRepo {
	return &WebUserRepo{db: db}
}

func (repo *WebUserRepo) GetWebUserByID(userID int) (*model.WebUser, error) {
	var user model.WebUser

	err := repo.db.Get(&user, "SELECT * FROM web_user WHERE id = $1;", userID)
	if err == sql.ErrNoRows {
		return nil, ErrWebUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting web user by id: %w", err)
	}

	return &user, nil
}

func (repo *WebUserRepo) GetWebUserByEmail(email string) (*model.WebUser, error) {
	var user model.WebUser

	err := repo.db.Get(&user, "SELECT * FROM web_user WHERE email = $1;", email)
	if err == sql.ErrNoRows {
		return nil, ErrWebUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error while getting web user by email: %w", err)
	}

	return &user, nil
}

func (repo *WebUserRepo) CreateWebUser(user *model.WebUser) (int, error) {
	var id int

	row := repo.db.QueryRow("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;", user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating web user: %w", err)
	}

	return id, nil
}
