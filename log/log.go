package log

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"gopkg.in/yaml.v2"

	"github.com/arutselvan15/go-utils/diff"
	"github.com/arutselvan15/go-utils/logconstants"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

//GoLog log
type GoLog interface {
	logrus.FieldLogger
	GetEntry() *logrus.Entry
	SetLevel(string) *Log
	SetCluster(string) *Log
	SetComponent(string) *Log
	SetSubComponent(string) *Log
	SetProcess(process string) *Log
	SetSubProcess(subProcess string) *Log
	SetAction(action string) *Log
	SetUser(user string) *Log
	SetInvolvedObj(involvedObj string) *Log
	SetDisposition(disposition string) *Log
	SetAPIRequest(string, string) *Log
	SetAPIResponse(string, string) *Log
	Event(string, ...interface{})
	SetObjectAudit(...interface{}) *Log
	GetLogger() *Log
	PushContext()
	PopContext()
	SaveContext()
	RestoreContext()
	PushPop(f func())
	ThreadLogger() *Log
}

// Log log
type Log struct {
	*logrus.Entry

	// creator
	logger *logrus.Logger

	// used for push/pop of contexts
	contextStack *stack.Stack

	// save/restore context stack
	savedContexts *stack.Stack

	// fields
	cluster         string
	component       string
	subComponent    string
	process         string
	subProcess      string
	action          string
	user            string
	involvedObj     string
	disposition     string
	apiEndpoint     string
	apiRequest      string
	apiResponse     string
	objectAuditType string
	objectAuditData string
}

// GetEntry get entry
func (l *Log) GetEntry() *logrus.Entry {
	return l.Entry
}

// SetObjectAudit sets the object audit data and type
func (l *Log) SetObjectAudit(objects ...interface{}) *Log {
	if len(objects) == 1 {
		// create
		l.objectAuditType = logconstants.Create
		l.Entry = l.WithField("objectAuditType", l.objectAuditType)
		yamlBytes, _ := yaml.Marshal(objects[0])
		l.objectAuditData = string(yamlBytes)
		l.Entry = l.WithField("objectAuditData", l.objectAuditData)

	} else if len(objects) == 2 {
		// update and look at first 2
		l.objectAuditType = logconstants.Update
		l.Entry = l.WithField("objectAuditType", l.objectAuditType)
		ch, _ := diff.GetDiffChangelog(objects[0], objects[1])
		data := ""
		if ch != nil {
			for _, c := range *ch {
				data = data + fmt.Sprintf("(%v, %v, %v, %v)\n", c.Path, c.Type, c.From, c.To)
			}
			l.objectAuditData = data
			l.Entry = l.WithField("objectAuditData", l.objectAuditData)
		}
	}

	return l
}

// SetLevel sets the level at which log messages are published/written.
func (l *Log) SetLevel(level string) *Log {
	// If there's no explicit logging level specified, set the level to INFO
	if level == "" {
		level = "info"
	}

	loglevel, err := logrus.ParseLevel(level)
	if err == nil {
		// set default logger and the custom logger levels
		logrus.SetLevel(loglevel)
		l.Logger.SetLevel(loglevel)
	}

	return l
}

// SetAPIRequest sets the endpoint, payload, and response of an api call
func (l *Log) SetAPIRequest(endpoint, request string) *Log {

	if endpoint == "" {
		delete(l.Data, "apiEndpoint")
		l.apiEndpoint = ""
	} else {
		l.apiEndpoint = endpoint
		l.Entry = l.WithField("apiEndpoint", l.apiEndpoint)
	}

	if request == "" {
		delete(l.Data, "apiRequest")
		l.apiRequest = ""
	} else {
		l.apiRequest = request
		l.Entry = l.WithField("apiRequest", l.apiRequest)
	}

	return l
}

// SetAPIResponse SetAPIResponse
func (l *Log) SetAPIResponse(endpoint, response string) *Log {

	if endpoint == "" {
		delete(l.Data, "apiEndpoint")
		l.apiEndpoint = ""
	} else {
		l.apiEndpoint = endpoint
		l.Entry = l.WithField("apiEndpoint", l.apiEndpoint)
	}

	if response == "" {
		delete(l.Data, "apiResponse")
		l.apiResponse = ""
	} else {
		l.apiResponse = response
		l.Entry = l.WithField("apiResponse", l.apiResponse)
	}

	return l
}

// SetCluster adds cluster name
func (l *Log) SetCluster(cluster string) *Log {
	if cluster == "" {
		delete(l.Data, "cluster")
		l.cluster = ""
	} else {
		l.cluster = cluster
		l.Entry = l.WithField("cluster", l.cluster)
	}

	return l
}

// SetComponent adds the component field to each Log message if provided
func (l *Log) SetComponent(c string) *Log {
	if c == "" {
		delete(l.Data, "component")
		l.component = ""
	} else {
		l.component = c
		l.Entry = l.WithField("component", l.component)
	}

	return l
}

// SetSubComponent adds the sub component field to each Log message if provided
func (l *Log) SetSubComponent(sc string) *Log {
	if sc == "" {
		delete(l.Data, "subComponent")
		l.subComponent = ""
	} else {
		l.subComponent = sc
		l.Entry = l.WithField("subComponent", l.subComponent)
	}

	return l
}

// SetProcess adds the process (CreateProject, DeleteProject, etc.) field to each Log message if provided
func (l *Log) SetProcess(process string) *Log {
	if process == "" {
		delete(l.Data, "process")
		l.process = ""
	} else {
		l.process = process
		l.Entry = l.WithField("process", l.process)
	}

	return l
}

// SetSubProcess adds the subProcess field to each log message if provided
func (l *Log) SetSubProcess(subProcess string) *Log {
	if subProcess == "" {
		delete(l.Data, "subProcess")
		l.subProcess = ""
	} else {
		l.subProcess = subProcess
		l.Entry = l.WithField("subProcess", l.subProcess)
	}

	return l
}

// SetAction adds the action into the log (configureProjectAdmin, configureRBAC and etc.)
func (l *Log) SetAction(action string) *Log {
	if action == "" {
		delete(l.Data, "action")
		l.action = ""
	} else {
		l.action = action
		l.Entry = l.WithField("action", l.action)
	}

	return l
}

// SetUser adds the user field to each log message if provided
func (l *Log) SetUser(user string) *Log {
	if user == "" {
		delete(l.Data, "user")
		l.user = ""
	} else {
		l.user = user
		l.Entry = l.WithField("user", l.user)
	}

	return l
}

// SetInvolvedObj adds the involved object (route name, project name, etc.) field to each log message if provided
func (l *Log) SetInvolvedObj(involvedObj string) *Log {
	if involvedObj == "" {
		delete(l.Data, "involvedObject")
		l.involvedObj = ""
	} else {
		l.involvedObj = involvedObj
		l.Entry = l.WithField("involvedObject", l.involvedObj)
	}

	return l
}

// SetDisposition adds the disposition (Success/Fail of the process and/or action) field to each log message if provided
func (l *Log) SetDisposition(disposition string) *Log {
	if disposition == "" {
		delete(l.Data, "disposition")
		l.disposition = ""
	} else {
		l.disposition = disposition
		l.Entry = l.WithField("disposition", l.disposition)
	}

	return l
}

// Event Event
func (l *Log) Event(format string, args ...interface{}) {
	l.PushContext()
	l.SetAction("Event").Infof(format, args...)
	l.PopContext()
}

// SaveContext SaveContext
// Save the current context for later restore.
// This saves the current fields and the current stack.
// After this call, the current context is intact, but has a new empty stack.
func (l *Log) SaveContext() {
	// create a new context
	c := l.GetLogger()
	// copy state from existing
	c.copyContextFrom(l)
	// save the current context stack state
	c.contextStack = l.contextStack
	// Push current context for later restore
	l.savedContexts.Push(c)

	// Create new context stack for new scope
	l.contextStack = stack.New()
	// Now, current logger is unchanged except for a new push/pop stack, and the current whole context (including
	// stack) is saved in the saved stack.
}

// RestoreContext function
// If previously saved, this restores the saved context.
// This restores the previouosly saved field and the previous stack.
// Anything in the stack between Save/Restore is gone.
func (l *Log) RestoreContext() {
	// If no saves, do nothing
	if l.savedContexts.Len() == 0 {
		return
	}
	// Get the last saved full context
	c := l.savedContexts.Pop().(*Log)
	// Restore saved push/pop stack
	l.contextStack = c.contextStack
	// Restore context data
	l.copyContextFrom(c)
}

// PushContext PushContext
func (l *Log) PushContext() {
	// push and pop by value, not by reference
	l.contextStack.Push(*l)
}

// PopContext PopContext
func (l *Log) PopContext() {
	// Do nothing if nothing there
	if l.contextStack.Len() == 0 {
		return
	}

	pop := l.contextStack.Pop().(Log)

	l.copyContextFrom(&pop)
}

// PushPop function
// Run the function within a push/pop
func (l *Log) PushPop(f func()) {
	l.PushContext()
	f()
	l.PopContext()
}

// Doesn't copy the stack, just the fields
func (l *Log) copyContextFrom(from *Log) *Log {
	l.logger = from.logger
	l.Entry = from.logger.WithFields(logrus.Fields{})
	l.SetComponent(from.component)
	l.SetCluster(from.cluster)
	l.SetSubComponent(from.subComponent)
	l.SetProcess(from.process)
	l.SetSubProcess(from.subProcess)
	l.SetAction(from.action)
	l.SetUser(from.user)
	l.SetInvolvedObj(from.involvedObj)
	l.SetDisposition(from.disposition)
	l.SetAPIRequest(from.apiEndpoint, from.apiRequest)
	l.SetAPIResponse(from.apiEndpoint, from.apiResponse)
	l.setObjectAudit(l.objectAuditType, l.objectAuditData)

	return l
}

func (l *Log) setObjectAudit(auditType string, data string) {
	if auditType != "" {
		l.objectAuditType = auditType
		l.Entry = l.WithField("objectAuditType", l.objectAuditType)
	}

	if data != "" {
		l.objectAuditData = data
		l.Entry = l.WithField("objectAuditData", l.objectAuditData)
	}
}

func (l *Log) clear() {
	l.Entry = l.logger.WithFields(logrus.Fields{})

	l.cluster = ""
	l.component = ""
	l.subComponent = ""
	l.process = ""
	l.subProcess = ""
	l.action = ""
	l.user = ""
	l.involvedObj = ""
	l.disposition = ""
	l.apiEndpoint = ""
	l.apiRequest = ""
	l.apiResponse = ""
	l.objectAuditData = ""
	l.objectAuditType = ""
}

// GetLogger creates and returns a new logging context.
// A logging context is a wrapper for a log entry for a logger.  These objects can NOT be used
// in parallel.
// GetLogger GetLogger
func (l *Log) GetLogger() *Log {
	nl := newLog(l.logger, l.contextStack, l.savedContexts)
	return nl
}

// ThreadLogger clones the logger, but gives it a new set of stacks.
// This allows the thread to do all context stuff independently of other threads
// ThreadLogger ThreadLogger
func (l *Log) ThreadLogger() *Log {
	// create the context stacks
	stack1 := stack.New()
	stack2 := stack.New()

	// create the initial logging context
	nlog := newLog(l.logger, stack1, stack2)
	// populate with existing context
	nlog.copyContextFrom(l)
	return nlog
}

//NewLoggerWithFile log
func NewLoggerWithFile(filename string) GoLog {
	logger := logrus.New()

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  RFC3339NanoFixed,
	})

	// log file
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   filename,
		MaxSize:    250,
		MaxBackups: 3,
		MaxAge:     30,
		Level:      logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: RFC3339NanoFixed,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logger.AddHook(rotateFileHook)

	// create the shared context stacks
	stack1 := stack.New()
	stack2 := stack.New()

	// create the initial logging context
	return newLog(logger, stack1, stack2)
}

// NewLogger creates a logger context for a newly created logging output.
func NewLogger() GoLog {
	logger := logrus.New()

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  RFC3339NanoFixed,
	})

	// create the shared context stacks
	stack1 := stack.New()
	stack2 := stack.New()

	// create the initial logging context
	return newLog(logger, stack1, stack2)
}

func newLog(logger *logrus.Logger, contextStack *stack.Stack, savedContexts *stack.Stack) *Log {
	nl := &Log{
		logger:        logger,
		contextStack:  contextStack,
		savedContexts: savedContexts,
	}

	nl.clear()
	return nl
}
