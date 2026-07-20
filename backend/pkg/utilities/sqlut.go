package utilities

import (
	"strings"
	"time"
)

// placeholdersForInts builds a "?, ?, ?" placeholder string and matching
// []interface{} args for use in an `IN (...)` clause.
func PlaceholdersForInts(ids []int) (string, []interface{}) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	return strings.Join(placeholders, ", "), args
}

func ParseSQLiteTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}
