package main

import "log"

var Version = "DEV"

func main() {
	log.Printf("v%s\n", Version)
	Execute()
}
