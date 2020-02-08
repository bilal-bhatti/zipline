package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tcases = []struct {
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
	for _, tcase := range tcases {
		tt := newTypeToken(tcase.pkg, tcase.sig, "")

		assert.Equal(t, tcase.vn, tt.varName(), "Name should be same")
		assert.Equal(t, tcase.vnp, tt.varNameAsPointer(), "Name as pointer should be same")
		assert.Equal(t, tcase.inst, tt.inst(), "Instantiate should be same")
		assert.Equal(t, tcase.pkg, tt.pkg(), "Package should be same")
		assert.Equal(t, tcase.param, tt.param(), "Param name should be same")
	}
}
