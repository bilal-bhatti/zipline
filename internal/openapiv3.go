package internal

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

// quick and dirty conversion from v2 specs to v3
// TODO: switch to using newer v2 spec and v3 spec APIs
// getkin/kin-openapi seems to better support for v2/v3 going forward
func convertToV3() error {
	input, err := os.ReadFile("api.oasv2.json")
	if err != nil {
		return err
	}

	var v2 openapi2.T
	if err = json.Unmarshal(input, &v2); err != nil {
		return err
	}

	v3, err := openapi2conv.ToV3(&v2)
	if err != nil {
		return err
	}

	bites, err := json.MarshalIndent(v3, "", "  ")
	if err != nil {
		return err
	}

	bites = append(bites, []byte("\n")...)

	err = os.WriteFile("api.oasv3.json", bites, 0644)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Printf("wrote OpenAPI v3 spec to  %s\n", path.Join(cwd, "api.oasv3.json"))

	return nil
}
