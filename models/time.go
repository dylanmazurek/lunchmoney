package models

import (
	"encoding/json"
	"time"
)

type DateTimeOptional time.Time

func (dto *DateTimeOptional) UnmarshalJSON(data []byte) error {
	var err error
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	_, err = time.Parse(time.RFC3339, string(data))
	if err == nil { // is date and time
		return json.Unmarshal(data, &dto)
	}

	var dateString string
	err = json.Unmarshal(data, &dateString)
	if err != nil {
		return err
	}

	dateParsed, err := time.Parse(time.DateOnly, dateString)
	if err == nil { // is date only
		*dto = DateTimeOptional(dateParsed)
		return nil
	}

	return err
}

func (dto *DateTimeOptional) MarshalJSON() ([]byte, error) {
	return json.Marshal(dto)
}

func (dto *DateTimeOptional) Time() time.Time {
	return time.Time(*dto)
}

func (dto *DateTimeOptional) String() string {
	return dto.Time().String()
}
