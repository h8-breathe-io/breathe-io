package util

import "time"

func ParseStrToUnixDate(date string) (int64, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func ParseUnixToDate(unix int64) string {
	t := time.Unix(unix, 0)
	return t.Format("2006-01-02")
}
