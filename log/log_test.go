package log

import (
	"testing"
)

func TestLog_GetLogger(t *testing.T) {
	logger := NewLogger("myapp", FormatText, "").SetLevel(DebugLevel)
	logger.GetLogger().SetSubComponent("user").SetAction("create").SetStep("step1").Debug("step1 log")
}

func TestLog_GetLoggerWithFile(t *testing.T) {
	logger := NewLogger("myapp", FormatJson, "local.log").SetLevel(DebugLevel)
	logger.GetLogger().SetSubComponent("user").SetAction("create").SetStep("step1").Debug("step1 log with logfile")
}

func TestLog_GetLoggerAudit(t *testing.T) {
	logger := NewLogger("myapp", FormatText, "").SetLevel(DebugLevel)
	logger.GetLogger().SetSubComponent("user").SetAction("create").SetAudit("audit", "user1", "auditenabled", "auditdisabled").Debug("update performed")
}
