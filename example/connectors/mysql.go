package connectors

import (
	"fmt"
)

type DataConnector struct {
	Data map[int]string
}

var data = map[int]string{
	1: "one",
	2: "two",
}

// For a singleton DB connector
var dc = &DataConnector{
	Data: data,
}

func ProvideDataConnector() (*DataConnector, error) {
	return dc, nil
}

func (dc *DataConnector) Create(input string) string {
	dc.Data[len(dc.Data)+1] = input
	return fmt.Sprint(len(dc.Data))
}

func (dc *DataConnector) GetOne(id int) string {
	return dc.Data[id]
}
