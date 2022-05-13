package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type WebUserRepo struct {
	db *DB
}

func NewWebUserRepo(db *DB) *WebUserRepo {
	return &WebUserRepo{db: db}
}

func (repo *WebUserRepo) GetWebUser(userID int) (*model.WebUser, error) {
	user := &model.WebUser{}

	rows, err := repo.db.Query("SELECT * FROM web_user WHERE id=$1;", userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting web user: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			continue
		}
	}

	if user.Email == "" || user.Password == "" {
		return nil, nil
	}

	return user, nil
}

func (repo *WebUserRepo) GetWebUserByEmail(email string) (*model.WebUser, error) {
	user := &model.WebUser{}

	rows, err := repo.db.Query("SELECT * FROM web_user WHERE email=$1;", email)
	if err != nil {
		return nil, fmt.Errorf("error while getting web user by email: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			continue
		}
	}

	if user.Email == "" || user.Password == "" {
		return nil, nil
	}

	return user, nil
}

func (repo *WebUserRepo) CreateWebUser(user *model.WebUser) (int, error) {
	var id int

	row := repo.db.QueryRow("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;", user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating web user: %w", err)
	}

	return id, nil
}
