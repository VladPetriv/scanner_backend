package mocks

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreateMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, fmt.Errorf("create mock: %w", err)
	}

	return db, mock, nil
}
