package web

import (
	"strconv"
)

func ParseInt64(strings []string) ([]int64, error) {
	var ints = []int64{}

	for _, i := range strings {
		j, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return ints, err
		}
		ints = append(ints, j)
	}

	return ints, nil
}
