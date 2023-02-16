package pg

import (
	"database/sql"
	"errors"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type WebUserRepo struct {
	db *DB
}

func NewWebUserRepo(db *DB) *WebUserRepo {
	return &WebUserRepo{db: db}
}

func (repo WebUserRepo) CreateWebUser(user *model.WebUser) error {
	_, err := repo.db.Exec(
		"INSERT INTO web_user(email, password) VALUES ($1, $2);",
		user.Email, user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo WebUserRepo) GetWebUserByEmail(email string) (*model.WebUser, error) {
	var user model.WebUser

	err := repo.db.Get(&user, "SELECT * FROM web_user WHERE email = $1;", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
