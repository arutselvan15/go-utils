package log

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/arutselvan15/go-utils/logconstants"
	"github.com/golang-collections/collections/stack"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Setenv("CLUSTER", "test_cluster")
}

func logAndAssertJSON(_ *testing.T, log *Log, message string, assertions func(fields logrus.Fields)) {
	var (
		buffer bytes.Buffer
		fields logrus.Fields
	)

	log.Entry.Logger.Out = &buffer
	log.Entry.Logger.Formatter = new(logrus.JSONFormatter)
	log.Info(message)

	_ = json.Unmarshal(buffer.Bytes(), &fields)

	assertions(fields)
}

func newLogger() *Log {
	return newLog(logrus.New(), stack.New(), stack.New())
}

func TestGetLogger(t *testing.T) {
	logger := newLogger()
	logger.GetLogger()
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})
}

func TestSetObjectName(t *testing.T) {
	involvedObj := "test_involvedObject"
	logger := newLogger()
	logger.SetObjectName(involvedObj)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, involvedObj, fields["objectName"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetObjectName("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["objectName"])
	})
}

func TestSetUser(t *testing.T) {
	user := "test_user"
	logger := newLogger()
	logger.SetUser(user)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, user, fields["user"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetUser("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["user"])
	})
}

func TestSetOperation(t *testing.T) {
	action := "operation"
	logger := newLogger()
	logger.SetOperation(action)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, action, fields["operation"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetOperation("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["operation"])
	})
}

func TestSetLevel(t *testing.T) {
	logger := newLogger()
	logger.SetLevel(InfoLevel)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	// Log event is of type info will not be processed because Log level should be
	// panic and above
	logger.SetLevel(PanicLevel)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, nil, fields["msg"])
		assert.Equal(t, nil, fields["level"])
	})
}

func TestAudits(t *testing.T) {
	logger := newLogger()
	logger.LogAuditEvent("")
	logger.LogAuditObject([]string{})
	logger.LogAuditObject([]string{}, []string{})
	logger.LogAuditAPI("GET", "/", "", "", 0)
}

func TestFormatters(t *testing.T) {
	logger := newLogger()
	flogger := NewLoggerWithFile("", 1, 1, 1)

	logger.SetFormatterType(TextFormatterType)
	logger.SetFormatterType(JSONFormatterType)

	flogger.SetLogFileFormatterType(TextFormatterType)
	flogger.SetLogFileFormatterType(JSONFormatterType)
}

func TestPushPop(t *testing.T) {
	logger := newLogger()
	logger.PushPop(func() {})
}

func assertAllFields(t *testing.T, log *Log, prefix string) {
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, prefix+"cluster", fields["cluster"])
		assert.Equal(t, prefix+"app", fields["app"])
		assert.Equal(t, prefix+"resource", fields["resource"])
		assert.Equal(t, prefix+"component", fields["component"])
		assert.Equal(t, prefix+"operation", fields["operation"])
		assert.Equal(t, prefix+"objectname", fields["objectName"])
		assert.Equal(t, prefix+"objectstate", fields["objectState"])
		assert.Equal(t, prefix+"user", fields["user"])
		assert.Equal(t, prefix+"step", fields["step"])
		assert.Equal(t, prefix+"stepstate", fields["stepState"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})
}

func setAllFields(log *Log, prefix string) {
	log.SetCluster(prefix + "cluster")
	log.SetApplication(prefix + "app")
	log.SetResource(prefix + "resource")
	log.SetComponent(prefix + "component")
	log.SetOperation(prefix + "operation")
	log.SetObjectName(prefix + "objectname")
	log.SetObjectState(prefix + "objectstate")
	log.SetUser(prefix + "user")
	log.SetStep(prefix + "step")
	log.SetStepState(prefix + "stepstate")
	log.SetLevel("info")
}

func TestAllFields(t *testing.T) {
	logger := newLogger()
	setAllFields(logger, "main")
	assertAllFields(t, logger, "main")
}

func TestMultiLog(t *testing.T) {
	log1 := newLogger()
	log2 := log1.GetLogger()
	assert.Equal(t, log1, log2, "Initial logs should be equal")

	log2.SetOperation("test")
	assert.NotEqual(t, log1, log2, "Changed logs should not be equal")

	logAndAssertJSON(t, log1, "test", func(fields logrus.Fields) {
		assert.Equal(t, nil, fields["operation"])
	})
	logAndAssertJSON(t, log2, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["operation"])
	})
	log1.SetOperation("hahaha")
	logAndAssertJSON(t, log1, "blahblah", func(fields logrus.Fields) {
		assert.Equal(t, "hahaha", fields["operation"])
	})
	logAndAssertJSON(t, log2, "yeet", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["operation"])
	})
}

func TestLogStack(t *testing.T) {
	log := newLogger()

	log.SetComponent("TestLogStack").SetOperation("BeforePush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["component"])
		assert.Equal(t, "BeforePush", fields["operation"])
	})

	log.PushContext()
	log.SetObjectName("testobject")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["component"])
		assert.Equal(t, "BeforePush", fields["operation"])
		assert.Equal(t, "testobject", fields["objectName"])
	})

	log.SetOperation("AfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["component"])
		assert.Equal(t, "AfterPush", fields["operation"])
		assert.Equal(t, "testobject", fields["objectName"])
	})

	log.PopContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["component"])
		assert.Equal(t, "BeforePush", fields["operation"])
		assert.Nil(t, fields["objectName"]) // pop should clear the nested set
	})
}

func saveRestore(log *Log, t *testing.T) {
	log.SetComponent("TestSaveRestore").SetOperation("BeforeSave")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSave", fields["operation"])
	})

	log.PushContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSave", fields["operation"])
	})

	log.SetOperation("BeforeSaveAfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["operation"])
	})

	log.SaveContext()

	// New stack should be empty
	assert.Equal(t, log.contextStack.Len(), 0)

	// Context should still be there
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["operation"])
	})
	log.SetOperation("AfterSave")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "AfterSave", fields["operation"])
	})

	// Pushing within new context stack
	log.PushContext()
	log.PushContext()
	assert.Equal(t, log.contextStack.Len(), 2)

	log.SetOperation("AfterSaveAfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "AfterSaveAfterPush", fields["operation"])
	})

	// Now restore the context - new stack should be replaced by old one
	log.RestoreContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["operation"])
	})
	assert.Equal(t, log.contextStack.Len(), 1)

	// Back to initial
	log.PopContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["component"])
		assert.Equal(t, "BeforeSave", fields["operation"])
	})
}

func TestSaveRestore(t *testing.T) {
	log := newLogger()
	saveRestore(log, t)
}

func TestThreadedContexts(t *testing.T) {
	log := newLogger()

	setAllFields(log, "initial")
	assertAllFields(t, log, "initial")

	t1log := log.ThreadLogger()
	t2log := log.ThreadLogger()

	// Ensure both are cloned fully
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "initial")

	// Thread 1
	t1log.PushContext()
	assertAllFields(t, t1log, "initial")
	// change all fields, and check
	setAllFields(t1log, "t1log")
	assertAllFields(t, t1log, "t1log")
	assertAllFields(t, t2log, "initial")

	// Thread 2
	t2log.PushContext()
	// change all fields in t2, and check
	setAllFields(t2log, "t2log")
	assertAllFields(t, t2log, "t2log")
	// make sure t1log is ok...
	assertAllFields(t, t1log, "t1log")

	// Thead 1
	// Pop t1 context, and make sure t2 is not affected...
	t1log.PopContext()
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "t2log")
	// Make usre t2log's stack wasn't popped
	assert.Equal(t, 1, t2log.contextStack.Len())

	// Thread 2
	t2log.PopContext()
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "initial")
}

func TestNestedSaveRestore(t *testing.T) {
	log := newLogger()

	log.SetComponent("TestNestedSaveRestore").SetOperation("Initial")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestNestedSaveRestore", fields["component"])
		assert.Equal(t, "Initial", fields["operation"])
	})

	log.PushContext()
	log.PushContext()
	log.PushContext()
	assert.Equal(t, log.contextStack.Len(), 3)

	// Top level save
	log.SaveContext()

	// After save, the context stack is empty, and the saved stack has 1
	assert.Equal(t, log.contextStack.Len(), 0)
	assert.Equal(t, log.savedContexts.Len(), 1)

	// call nested save restore testing
	saveRestore(log, t)

	// Restore to initial
	log.RestoreContext()

	// After restore, the initial context stack is restored, and the saved stack is empty
	assert.Equal(t, log.contextStack.Len(), 3)
	assert.Equal(t, log.savedContexts.Len(), 0)

	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestNestedSaveRestore", fields["component"])
		assert.Equal(t, "Initial", fields["operation"])
	})
}

func TestThreadedSaveRestore(t *testing.T) {
	log := newLogger()

	setAllFields(log, "initial")
	assertAllFields(t, log, "initial")

	t1log := log.ThreadLogger()
	t2log := log.ThreadLogger()

	// Ensure both are cloned fully
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "initial")

	// Thread 1
	t1log.SaveContext()
	assertAllFields(t, t1log, "initial")
	// change all fields, and check
	setAllFields(t1log, "t1log")
	assertAllFields(t, t1log, "t1log")
	assertAllFields(t, t2log, "initial")

	// Thread 2
	t2log.SaveContext()
	// change all fields in t2, and check
	setAllFields(t2log, "t2log")
	assertAllFields(t, t2log, "t2log")
	// make sure t1log is ok...
	assertAllFields(t, t1log, "t1log")

	// Thead 1
	// Pop t1 context, and make sure t2 is not affected...
	t1log.RestoreContext()
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "t2log")
	// Make usre t2log's stack wasn't popped
	assert.Equal(t, 1, t2log.savedContexts.Len())

	// Thread 2
	t2log.RestoreContext()
	assertAllFields(t, t1log, "initial")
	assertAllFields(t, t2log, "initial")
}

func TestAllLogConstants(t *testing.T) {
	logger := newLogger()

	actions := []string{
		logconstants.Create, logconstants.Update, logconstants.Delete, logconstants.Read,
		logconstants.Mutate, logconstants.Audit, logconstants.Validate,
		logconstants.Pending, logconstants.InProgress, logconstants.Error, logconstants.Retry, logconstants.Unknown, logconstants.Deleting,
		logconstants.Successful, logconstants.Failed, logconstants.Processing,
	}

	for _, action := range actions {
		logger.SetOperation(action)
		logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) { assert.Equal(t, action, fields["operation"]) })
	}
}

func TestNewLoggerWithFile(t *testing.T) {
	type args struct {
		component string
		filename  string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success new logger with file",
			args: args{component: "test_component", filename: "/tmp/test.log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLoggerWithFile(tt.args.filename, 1, 1, 1)
			if got == nil {
				t.Errorf("NewLoggerWithFile() = %v, want customlog", got)
				return
			}
			got.Info("test log")
		})
	}
}

func TestNewLogger(t *testing.T) {
	type args struct {
		component string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success newlogger",
			args: args{component: "test_component"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogger()
			if got == nil {
				t.Errorf("NewLogger() = %v, want customlog", got)
				return
			}
			got.Info("test log")
		})
	}
}
