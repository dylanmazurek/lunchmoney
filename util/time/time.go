package time

import (
	"encoding/json"
	"strconv"
	"time"
)

type DateTime struct {
	Time time.Time
}

type Date struct {
	Date time.Time
}

func (t *Date) String() string {
	return t.Date.Format(time.DateOnly)
}

func (t *Date) MarshalJSON() ([]byte, error) {
	type Alias Date
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})

	return marshaledJSON, err
}

func (t *Date) UnmarshalJSON(data []byte) error {
	dateString := string(data)
	if dateString == "null" {
		return nil
	}

	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	t.Date, err = time.Parse(time.DateOnly, s)

	return err
}

func Parse(date string) Date {
	t, _ := time.Parse(time.DateOnly, date)
	return Date{Date: t}
}
