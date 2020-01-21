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
	l.Info("sample stdout log")

	l = log.NewLoggerWithFile("/tmp/tlog.log", 1, 1, 1)
	l.Info("sample file log")
	l.SetLogFileFormatterType(log.JsonFormatterType)
	l.Info("sample file log json format")
}
