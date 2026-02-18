package shared

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// TimeToDbDate transforms a given date string to a pgtype.Date
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

// ParseDate parses a given date string to a time.Time object
func ParseDate(date string) (time.Time, error) {
	format := "2006-01-02"
	return time.Parse(format, date)

}
