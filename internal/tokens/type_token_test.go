package tokens

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tCase struct {
	pkg,
	fqSignature,
	varName,
	newInstance,
	declSignature,
	importingPkg string
}

var tcases = []tCase{
	{
		fqSignature:   "github.com/bilal-bhatti/zipline/example/web.ContactResponse",
		varName:       "contactResponse",
		newInstance:   "ContactResponse{}",
		declSignature: "ContactResponse",
		importingPkg:  "github.com/bilal-bhatti/zipline/example/web",
	},
	{
		fqSignature:   "*github.com/bilal-bhatti/zipline/example/web.ContactResponse",
		varName:       "contactResponse",
		newInstance:   "&ContactResponse{}",
		declSignature: "*ContactResponse",
		importingPkg:  "github.com/bilal-bhatti/zipline/example/web",
	},
	{
		fqSignature:   "*github.com/bilal-bhatti/zipline/example/models.ContactResponse",
		varName:       "contactResponse",
		newInstance:   "&models.ContactResponse{}",
		declSignature: "*models.ContactResponse",
		importingPkg:  "github.com/bilal-bhatti/zipline/example/web",
	},
	{
		fqSignature:   "github.com/bilal-bhatti/zipline/example/models.ContactResponse",
		varName:       "contactResponse",
		newInstance:   "models.ContactResponse{}",
		declSignature: "models.ContactResponse",
		importingPkg:  "github.com/bilal-bhatti/zipline/example/web",
	},
}

func TestVarTokenParse(t *testing.T) {
	for _, tcase := range tcases {
		tt := NewTypeToken(tcase.declSignature, "")

		assert.Equal(t, tcase.varName, tt.VarName(), "Name should be same")
		assert.Equal(t, tcase.newInstance, tt.NewInstance(tcase.importingPkg), "Instantiate should be same")
		assert.Equal(t, tcase.declSignature, tt.DeclSignature(tcase.importingPkg), "Param name should be same")
	}
}
