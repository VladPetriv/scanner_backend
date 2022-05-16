package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type SavedRepo struct {
	db *DB
}

func NewSavedRepo(db *DB) *SavedRepo {
	return &SavedRepo{db: db}
}

func (repo *SavedRepo) GetSavedMessages(UserID int) ([]model.Saved, error) {
	savedMessages := make([]model.Saved, 0)

	rows, err := repo.db.Query("SELECT * FROM saved WHERE user_id=$1;", UserID)
	if err != nil {
		return nil, fmt.Errorf("error while getting saved messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		savedMessage := model.Saved{}
		err := rows.Scan(&savedMessage.ID, &savedMessage.WebUserID, &savedMessage.MessageID)
		if err != nil {
			continue
		}

		savedMessages = append(savedMessages, savedMessage)
	}

	if len(savedMessages) == 0 {
		return nil, nil
	}

	return savedMessages, nil
}

func (repo *SavedRepo) GetSavedMessageByMessageID(ID int) (*model.Saved, error) {
	savedMessage := &model.Saved{}

	rows, err := repo.db.Query("SELECT * FROM saved WHERE message_id=$1;", ID)
	if err != nil {
		return nil, fmt.Errorf("error while getting saved message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&savedMessage.ID, &savedMessage.WebUserID, &savedMessage.MessageID)
		if err != nil {
			continue
		}
	}

	if savedMessage.WebUserID == 0 || savedMessage.MessageID == 0 {
		return nil, nil
	}

	return savedMessage, nil
}

func (repo *SavedRepo) CreateSavedMessage(saved *model.Saved) (int, error) {
	var id int

	row := repo.db.QueryRow("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;", saved.WebUserID, saved.MessageID)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating saved message: %w", err)
	}

	return 1, nil
}

func (repo *SavedRepo) DeleteSavedMessage(ID int) (int, error) {
	var id int

	row := repo.db.QueryRow("DELETE FROM saved WHERE id=$1 RETURNING id;", ID)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while deleting saved message: %w", err)
	}

	return 1, nil
}
