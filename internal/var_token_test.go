package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var vts = []struct {
	sig, vn, vnp, inst, pkg, param string
}{
	{
		"ContactResponse",
		"contactResponse",
		"&contactResponse",
		"ContactResponse{}",
		"",
		"ContactResponse",
	},
	{
		"*ContactResponse",
		"contactResponse",
		"contactResponse",
		"&ContactResponse{}",
		"",
		"*ContactResponse",
	},
	{
		"*github.com/bilal-bhatti/zipline/example/models.ContactResponse",
		"contactResponse",
		"contactResponse",
		"&models.ContactResponse{}",
		"github.com/bilal-bhatti/zipline/example/models",
		"*models.ContactResponse",
	},
	{
		"github.com/bilal-bhatti/zipline/example/models.ContactResponse",
		"contactResponse",
		"&contactResponse",
		"models.ContactResponse{}",
		"github.com/bilal-bhatti/zipline/example/models",
		"models.ContactResponse",
	},
	{
		"*zipline/example/models.ContactResponse",
		"contactResponse",
		"contactResponse",
		"&models.ContactResponse{}",
		"zipline/example/models",
		"*models.ContactResponse",
	},
	{
		"zipline/example/models.ContactResponse",
		"contactResponse",
		"&contactResponse",
		"models.ContactResponse{}",
		"zipline/example/models",
		"models.ContactResponse",
	},
}

func TestVarTokenParse(t *testing.T) {
	for _, va := range vts {
		vt := newVarToken(va.pkg, va.sig, "")

		assert.Equal(t, va.vn, vt.varName(), "Name should be same")
		assert.Equal(t, va.vnp, vt.varNameAsPointer(), "Name as pointer should be same")
		assert.Equal(t, va.inst, vt.inst(), "Instantiate should be same")
		assert.Equal(t, va.pkg, vt.pkg(), "Package should be same")
		assert.Equal(t, va.param, vt.param(), "Param name should be same")
	}
}
