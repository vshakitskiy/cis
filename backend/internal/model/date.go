package model

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(dateFormat) + `"`), nil
}

// Tells pgx how to send Date to PostgreSQL
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

// Tells pgx ow to read a PostgreSQL date back into Date
func (d *Date) Scan(src any) error {
	if t, ok := src.(time.Time); ok {
		d.Time = t
		return nil
	}
	return errors.New("unsupported type for Date")
}
