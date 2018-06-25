package lib

import "time"

const httpTimePattern = "2006-01-02 15:04:05.0"

func toHTTPTimeStr(t time.Time) string {
	return t.Format(httpTimePattern)
}
func fromHTTPTimeStr(s string) (time.Time, error) {
	return time.Parse(httpTimePattern, s)
}
