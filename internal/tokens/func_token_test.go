package tokens

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
		vt := FuncToken{
			Signature: functest.sig,
			Rets:      []*TypeToken{NewTypeToken("", "context.Context", "context")},
		}

		assert.Equal(t, functest.pkg, vt.Pkg(), "Package should be same")
		// assert.Equal(t, functest.call, vt.call(), "Call as pointer should be same")
	}
}
