package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var functests = []struct {
	sig, pkg, call string
}{
	{
		"github.com/bilal-bhatti/zipline/example/services.ProvideContext",
		"github.com/bilal-bhatti/zipline/example/services",
		"services.ProvideContext(req)",
	},
}

func TestFuncTokenParse(t *testing.T) {
	for _, functest := range functests {
		vt := funcToken{
			signature: functest.sig,
		}

		assert.Equal(t, functest.pkg, vt.pkg(), "Package should be same")
		assert.Equal(t, functest.call, vt.call("req"), "Call as pointer should be same")
	}
}
