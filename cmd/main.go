package main

import (
	"flag"
	"github.com/arutselvan15/go-utils/log"
)

func main() {
	var action string

	flag.StringVar(&action, "action", action, "action name")
	flag.Parse()

	l := log.NewLogger()
	l.Info("sample log")
}
