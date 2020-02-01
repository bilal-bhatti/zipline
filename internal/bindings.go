package internal

type (
	binding struct {
		name, method, path string
		handler            *handlerInfo
	}

	handlerInfo struct {
		id      string
		params  []*varToken
		returns []*varToken
	}
)

func (b binding) id() string {
	return b.handler.id
}
