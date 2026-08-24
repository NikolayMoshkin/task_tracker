package helper

import "time"

func FormatOptionalTime(t *time.Time, layout string) string {
	if t == nil {
		return "-"
	}
	return t.Format(layout)
}
