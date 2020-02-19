package main

import (
	"fmt"
	"github.com/bilal-bhatti/zipline/internal"
)

func main() {
	zipline := internal.NewZipline()

	err := zipline.Start()

	if err != nil {
		fmt.Println(fmt.Errorf("zipline: %s", err.Error()))
	}
}
