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
		"context := services.ProvideContext()",
	},
}

func TestFuncTokenParse(t *testing.T) {
	for _, functest := range functests {
		vt := funcToken{
			signature: functest.sig,
			rets:      []*varToken{newVarToken("", "context.Context", "context")},
		}

		assert.Equal(t, functest.pkg, vt.pkg(), "Package should be same")
		assert.Equal(t, functest.call, vt.call(), "Call as pointer should be same")
	}
}
