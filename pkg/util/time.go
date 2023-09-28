package util

import "time"

func MustParseRFC3339(in string) time.Time {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		panic(err)
	}
	return t
}
