package web

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

func ParseTime(ts, format string) (*time.Time, error) {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, errors.Wrap(err, "invalid time zone")
	}

	var fmt = time.RFC3339
	if format != "" {
		parts := strings.Split(fmt, ",")
		// format should be of form `date-time,2006-01-02...` or `date-time`
		// this parsing function can be improved as needed.
		if len(parts) == 2 {
			fmt = parts[1]
		}
	}

	dt, err := time.ParseInLocation(fmt, ts, loc)

	if err != nil {
		return nil, err
	}

	return &dt, nil
}
