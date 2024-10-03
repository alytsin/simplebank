package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

var ErrNoRows = errors.New("no rows in result set")
var ErrUniqueViolation = errors.New("unique violation")

func TranslateError(err error) error {

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	var e *pq.Error
	if errors.As(err, &e) {
		switch e.Code {
		case "23505", "23503":
			// 23503 foreign key violation
			// 23505 unique violation
			return ErrUniqueViolation
		}
	}

	return err
}
