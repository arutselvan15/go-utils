package log

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/golang-collections/collections/stack"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/arutselvan15/go-utils/logconstants"
)

func init() {
	_ = os.Setenv("CLUSTER", "test_cluster")
}

func logAndAssertJSON(_ *testing.T, log *Log, message string, assertions func(fields logrus.Fields)) {
	var buffer bytes.Buffer
	var fields logrus.Fields

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

func TestSetDisposition(t *testing.T) {
	disposition := "test_disposition"
	logger := newLogger()
	logger.SetDisposition(disposition)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, disposition, fields["disposition"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetDisposition("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["disposition"])
	})
}

func TestSetInvolvedObj(t *testing.T) {
	involvedObj := "test_involvedObject"
	logger := newLogger()
	logger.SetInvolvedObj(involvedObj)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, involvedObj, fields["involvedObject"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetInvolvedObj("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["involvedObject"])
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

func TestSetAction(t *testing.T) {
	action := "test_action"
	logger := newLogger()
	logger.SetAction(action)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, action, fields["action"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetAction("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["action"])
	})
}

func TestSetProcess(t *testing.T) {
	process := "test_process"
	logger := newLogger()
	logger.SetProcess(process)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, process, fields["process"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetProcess("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["process"])
	})
}

func TestSetSubProcess(t *testing.T) {
	subProcess := "test_subprocess"
	logger := newLogger()
	logger.SetSubProcess(subProcess)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, subProcess, fields["subProcess"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetSubProcess("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["subProcess"])
	})
}

func TestSetSubComponent(t *testing.T) {
	subComponent := "test_sub_component"
	logger := newLogger()
	logger.SetSubComponent(subComponent)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, subComponent, fields["subComponent"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetSubComponent("")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["subComponent"])
	})
}

func TestSetApi(t *testing.T) {
	endpoint := "test_endpoint"
	request := "test_request"
	response := "test_response"
	logger := newLogger()

	logger.SetAPIRequest(endpoint, request)
	logger.SetAPIResponse(endpoint, response)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, endpoint, fields["apiEndpoint"])
		assert.Equal(t, request, fields["apiRequest"])
		assert.Equal(t, response, fields["apiResponse"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})

	logger.SetAPIRequest("", "")
	logger.SetAPIResponse("", "")
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Nil(t, fields["apiEndpoint"])
		assert.Nil(t, fields["apiRequest"])
		assert.Nil(t, fields["apiResponse"])
	})
}

func TestSetObjectAudit(t *testing.T) {
	type Quota struct {
		RAM     string
		Storage string
	}

	type ProjectSpec struct {
		Name          string
		Production    bool
		ProjectAdmins []string
		Quota         Quota
	}

	p1 := &ProjectSpec{Name: "test-p1", Production: true, ProjectAdmins: []string{"kemathew"}, Quota: Quota{RAM: "1Gi", Storage: "10Gi"}}
	p1v2 := &ProjectSpec{Name: "test-p1", Production: true, ProjectAdmins: []string{"kemathew"}, Quota: Quota{RAM: "1Gi", Storage: "11Gi"}}

	logger := newLogger()
	logger.SetObjectAudit(p1)
	yamlBytes, _ := yaml.Marshal(p1)

	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, logconstants.Create, fields["objectAuditType"])
		assert.Equal(t, string(yamlBytes), fields["objectAuditData"])
	})

	logger.SetObjectAudit(p1, p1v2)
	logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, logconstants.Update, fields["objectAuditType"])
		assert.Equal(t, "([Quota Storage], update, 10Gi, 11Gi)\n", fields["objectAuditData"])
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

func assertAllFields(t *testing.T, log *Log, prefix string) {
	endpoint := prefix + "test_endpoint"
	request := prefix + "test_request"
	response := prefix + "test_response"
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, endpoint, fields["apiEndpoint"])
		assert.Equal(t, request, fields["apiRequest"])
		assert.Equal(t, response, fields["apiResponse"])
		assert.Equal(t, prefix+"testAllFields", fields["process"])
		assert.Equal(t, prefix+"SUCCESS", fields["disposition"])
		assert.Equal(t, prefix+"Test", fields["action"])
		assert.Equal(t, prefix+"allFields", fields["subComponent"])
		assert.Equal(t, prefix+"unit-test", fields["involvedObject"])
		assert.Equal(t, prefix+"testuser", fields["user"])
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})
}

func setAllFields(log *Log, prefix string) {
	endpoint := prefix + "test_endpoint"
	request := prefix + "test_request"
	response := prefix + "test_response"
	log.SetAPIRequest(endpoint, request)
	log.SetAPIResponse(endpoint, response)
	log.SetProcess(prefix + "testAllFields")
	log.SetDisposition(prefix + "SUCCESS")
	log.SetAction(prefix + "Test")
	log.SetSubComponent(prefix + "allFields")
	log.SetInvolvedObj(prefix + "unit-test")
	log.SetLevel("info")
	log.SetUser(prefix + "testuser")
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

	log2.SetAction("test")
	assert.NotEqual(t, log1, log2, "Changed logs should not be equal")

	logAndAssertJSON(t, log1, "test", func(fields logrus.Fields) {
		assert.Equal(t, nil, fields["action"])
	})
	logAndAssertJSON(t, log2, "test", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["action"])
	})
	log1.SetAction("hahaha")
	logAndAssertJSON(t, log1, "blahblah", func(fields logrus.Fields) {
		assert.Equal(t, "hahaha", fields["action"])
	})
	logAndAssertJSON(t, log2, "yeet", func(fields logrus.Fields) {
		assert.Equal(t, "test", fields["action"])
	})
}

func TestLogStack(t *testing.T) {
	log := newLogger()

	log.SetProcess("TestLogStack").SetAction("BeforePush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["process"])
		assert.Equal(t, "BeforePush", fields["action"])
	})

	log.PushContext()
	log.SetInvolvedObj("testobject")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["process"])
		assert.Equal(t, "BeforePush", fields["action"])
		assert.Equal(t, "testobject", fields["involvedObject"])
	})

	log.SetAction("AfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["process"])
		assert.Equal(t, "AfterPush", fields["action"])
		assert.Equal(t, "testobject", fields["involvedObject"])
	})

	log.PopContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestLogStack", fields["process"])
		assert.Equal(t, "BeforePush", fields["action"])
		assert.Nil(t, fields["involvedObject"]) // pop should clear the nested set
	})
}

func saveRestore(log *Log, t *testing.T) {
	// Initial context
	log.SetProcess("TestSaveRestore").SetAction("BeforeSave")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSave", fields["action"])
	})

	log.PushContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSave", fields["action"])
	})

	// This will be the saved context
	log.SetAction("BeforeSaveAfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["action"])
	})

	log.SaveContext()

	// New stack should be empty
	assert.Equal(t, log.contextStack.Len(), 0)

	// Context should still be there
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["action"])
	})
	log.SetAction("AfterSave")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "AfterSave", fields["action"])
	})

	// Pushing within new context stack
	log.PushContext()
	log.PushContext()
	assert.Equal(t, log.contextStack.Len(), 2)

	log.SetAction("AfterSaveAfterPush")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "AfterSaveAfterPush", fields["action"])
	})

	// Now restore the context - new stack should be replaced by old one
	log.RestoreContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSaveAfterPush", fields["action"])
	})
	assert.Equal(t, log.contextStack.Len(), 1)

	// Back to initial
	log.PopContext()
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestSaveRestore", fields["process"])
		assert.Equal(t, "BeforeSave", fields["action"])
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

	log.SetProcess("TestNestedSaveRestore").SetAction("Initial")
	logAndAssertJSON(t, log, "test", func(fields logrus.Fields) {
		assert.Equal(t, "TestNestedSaveRestore", fields["process"])
		assert.Equal(t, "Initial", fields["action"])
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
		assert.Equal(t, "TestNestedSaveRestore", fields["process"])
		assert.Equal(t, "Initial", fields["action"])
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
		logconstants.OnAdd, logconstants.OnUpdate, logconstants.OnDelete,
		logconstants.Create, logconstants.Update, logconstants.Delete,
		logconstants.Mutate, logconstants.InProgress, logconstants.Read,
		logconstants.Audit, logconstants.Validate, logconstants.Retry, logconstants.Pending,
		logconstants.Start, logconstants.End, logconstants.Success, logconstants.Failure,
	}

	for _, action := range actions {
		logger.SetAction(action)
		logAndAssertJSON(t, logger, "test", func(fields logrus.Fields) { assert.Equal(t, action, fields["action"]) })
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
