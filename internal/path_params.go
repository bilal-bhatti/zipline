package internal

import (
	"fmt"
	"strings"
)

func isPathParam(binding *binding, vt *typeToken) bool {
	// TODO: fix to match regex patterns
	return strings.Contains(binding.path, fmt.Sprintf("{%s}", vt.name))
}
