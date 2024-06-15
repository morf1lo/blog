package exp

import "time"

func IsExpired(t time.Time) bool {
	now, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	return now.After(t)
}
