package utils

import (
	"net/url"
	"strings"
	"time"
)

func IsRobotsURL(u *url.URL) bool {
	const robotsTxtPath = "/robots.txt"

	if u == nil {
		return false
	}

	return strings.ToLower(u.Path) == robotsTxtPath
}

func IsTopic(u string) bool {
	return strings.Contains(u, "/viewtopic.php?")
}

func IsForum(u string) bool {
	return strings.Contains(u, "/viewforum.php?")
}

func ParseDate(date string) string {
	// 2017-08-26 22:40:42

	date = strings.TrimSpace(date)
	dateArr := strings.Split(date, " ")

	if dateArr[0] == "Today" || dateArr[0] == "Yesterday" {
		tNew := time.Now()

		if dateArr[0] == "Yesterday" {
			tNew = tNew.AddDate(0, 0, -1)
		}

		dateArr[0] = tNew.Format("2006-01-02")
		date = strings.Join(dateArr, " ")
	}

	t, err := time.Parse("2006-01-02 15:04:05", date)

	if err != nil {
		panic(err)
	}

	return t.UTC().Format(time.RFC3339)
}
