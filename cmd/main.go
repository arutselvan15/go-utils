package main

import (
	"flag"

	"github.com/arutselvan15/go-utils/testdata"
)

func main() {
	var action string

	flag.StringVar(&action, "action", action, "action name")
	flag.Parse()

	switch action {
	case "log":
		testdata.SampleLogging()
	default:
		testdata.SampleLogging()
	}
}
