package internal

import (
	"fmt"
	"strings"
)

func pathParam(binding *binding, vt *varToken) bool {
	return strings.Contains(binding.path, fmt.Sprintf("{%s}", vt.name))
}
