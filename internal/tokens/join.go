package tokens

import "strings"

func Join(tokens []*TypeToken, f func(t *TypeToken) string) string {
	s := []string{}
	for _, token := range tokens {
		s = append(s, f(token))
	}

	return strings.Join(s, ",")
}
