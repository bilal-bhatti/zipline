package internal

var (
	HREQ  = newTypeToken("", "*net/http.Request", "")
	HWRI  = newTypeToken("", "net/http.ResponseWriter", "")
	ERROR = newTypeToken("", "error", "err")
)

// global type key
const ZiplineTemplate = "ZiplineTemplate"
const ZiplineTemplateResolve = "Resolve"

const OpenAPIFile = "api.oasv2.json"
