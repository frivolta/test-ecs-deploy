package util

import (
	"time"
)

const dateLayout = "2006-01-02"

// ConvertDate Convert a string of type "YYYY-MM-DD" to sql readable time
func ConvertDate(str string) (time.Time, error) {
	conv, err := time.Parse(dateLayout, str)
	if err != nil {
		return time.Time{}, err
	}
	return conv, nil
}
