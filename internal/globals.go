package internal

var (
	HREQ  = newTypeToken("", "*net/http.Request", "")
	HWRI  = newTypeToken("", "net/http.ResponseWriter", "")
	ERROR = newTypeToken("", "error", "err")
)

// global type key
const (
	ZiplineTemplate        = "ZiplineTemplate"
	ZiplineTemplateResolve = "Resolve"
	ZiplineTemplateDevNull = "DevNull"
	ZiplineTemplateIgnore  = "Ignore"

	OpenAPIFile = "api.oasv2.json"
	Markdown    = "API.md"
)
