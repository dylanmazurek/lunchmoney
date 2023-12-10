package time

import (
	"encoding/json"
	"strconv"
	"time"
)

type DateTime struct {
	time.Time
}

type Date struct {
	time.Time
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

	t.Time, err = time.Parse(time.DateOnly, s)

	return err
}
