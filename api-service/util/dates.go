package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimeToDbDate(date string) (pgtype.Date, error) {
	format := "2006-01-02"

	parsedDate, err := time.Parse(format, date)

	if err != nil {
		return pgtype.Date{}, err
	}

	var dbDate pgtype.Date

	err = dbDate.Scan(parsedDate)

	if err != nil {
		return pgtype.Date{}, err
	}

	return dbDate, nil
}
