package utils

import "time"

func GetYMDFromDatabaseDate(dateString string) string {
	t, _ := time.Parse(time.RFC3339, dateString)
	// https://stackoverflow.com/a/20234207 WTF?
	return t.Format("2006-01-02")
}