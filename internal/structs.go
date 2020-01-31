package internal

type (
	binding struct {
		name, method, path string
		handler            *handlerInfo
	}

	handlerInfo struct {
		params  []*varToken
		returns []*varToken
	}
)

func (b binding) id() string {
	return "XYZ"
}
