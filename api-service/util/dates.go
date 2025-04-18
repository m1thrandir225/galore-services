package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Description:
// Given a date in the string format format it into a predefined format
// "YYYY-MM-DD" and return a pgtype.Date object from it
//
// Parameters:
// date: string
//
// Return:
// pgtype.Date
// error
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
