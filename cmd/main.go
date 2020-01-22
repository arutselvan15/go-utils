package main

import (
	"flag"

	"github.com/arutselvan15/go-utils/log"
	"github.com/arutselvan15/go-utils/logconstants"
)

func main() {
	var action string

	flag.StringVar(&action, "action", action, "action name")
	flag.Parse()

	l := log.NewLogger().SetComponent("logTest")

	l.SetAction(logconstants.Validate).SetState(logconstants.Start).Info("req received for validation")
	l.SetState(logconstants.InProgress).Info("validation in progress")
	l.SetState(logconstants.End).SetDisposition(logconstants.Success).Info("validation completed")

	l.SetAction("Provision").SetState(logconstants.Start).Info("req received for provision")
	l.SetState(logconstants.InProgress).Info("provision in progress")
	l.SetState(logconstants.End).SetDisposition(logconstants.Success).Info("provision completed")

	l.SetAction("Notification").SetState(logconstants.Start).Info("req received for notification")
	l.SetState(logconstants.InProgress).Info("notification in progress")
	l.SetState(logconstants.End).SetDisposition(logconstants.Failure).Info("notification failed")

	lf := log.NewLoggerWithFile("/tmp/tlog.log", 1, 1, 1)
	lf.Info("sample file log")
	lf.SetLogFileFormatterType(log.JsonFormatterType)
	lf.Info("sample file log json format")
}
