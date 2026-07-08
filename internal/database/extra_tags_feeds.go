package database

import "time"

func SqliteVarchar30ToTime(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

func TimeToSqliteVarchar30(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
}
