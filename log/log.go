// Package log lib
package log

import (
	"encoding/json"
	"fmt"

	"github.com/arutselvan15/go-utils/diff"
	"github.com/golang-collections/collections/stack"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// FormatterType formatter type
type FormatterType string

// LevelLog level log
type LevelLog string

var (
	// TextFormatterType TextFormatterType
	TextFormatterType FormatterType = "text"
	// JSONFormatterType JSONFormatterType
	JSONFormatterType FormatterType = "json"
	// DebugLevel DebugLevel
	DebugLevel LevelLog = "debug"
	// InfoLevel InfoLevel
	InfoLevel LevelLog = "info"
	// PanicLevel PanicLevel
	PanicLevel LevelLog = "panic"
)

//CommonLog log
type CommonLog interface {
	logrus.FieldLogger
	GetEntry() *logrus.Entry
	GetLogger() *Log
	SetLevel(LevelLog) *Log
	SetCluster(string) *Log
	SetApplication(string) *Log
	SetResource(string) *Log
	SetComponent(string) *Log
	SetOperation(action string) *Log
	SetObjectName(name string) *Log
	SetObjectState(state string) *Log
	SetUser(user string) *Log
	SetStep(string) *Log
	SetStepState(string) *Log
	LogAuditAPI(string, string, string, string, int)
	LogAuditObject(...interface{})
	LogAuditEvent(string)
	SetFormatterType(fType FormatterType) *Log
	SetLogFileFormatterType(fType FormatterType) *Log
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
	logger           *logrus.Logger
	logLevel         logrus.Level
	cluster          string
	application      string
	resource         string
	component        string
	operation        string
	objectName       string
	objectState      string
	user             string
	step             string
	stepState        string
	logFile          string
	logFileMaxSize   int
	logFileMaxAge    int
	logFileMaxBackup int

	// used for push/pop of contexts
	contextStack *stack.Stack

	// save/restore context stack
	savedContexts *stack.Stack
}

// GetEntry get entry
func (l *Log) GetEntry() *logrus.Entry {
	return l.Entry
}

// GetLogger creates and returns a new logging context.
// A logging context is a wrapper for a log entry for a logger.  These objects can NOT be used
// in parallel.
// GetLogger GetLogger
func (l *Log) GetLogger() *Log {
	nl := newLog(l.logger, l.contextStack, l.savedContexts)
	return nl
}

// SetLevel sets the level at which log messages are published/written.
func (l *Log) SetLevel(level LevelLog) *Log {
	loglevel, err := logrus.ParseLevel(string(level))
	if err != nil {
		// set default level on error
		loglevel, _ = logrus.ParseLevel(string(DebugLevel))
	}

	logrus.SetLevel(loglevel)
	l.Logger.SetLevel(loglevel)
	l.logLevel = loglevel

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

// SetApplication adds the app name
func (l *Log) SetApplication(app string) *Log {
	if app == "" {
		delete(l.Data, "app")
		l.application = ""
	} else {
		l.application = app
		l.Entry = l.WithField("app", l.application)
	}

	return l
}

// SetResource adds the resource
func (l *Log) SetResource(resource string) *Log {
	if resource == "" {
		delete(l.Data, "resource")
		l.resource = ""
	} else {
		l.resource = resource
		l.Entry = l.WithField("resource", l.resource)
	}

	return l
}

// SetComponent adds the component (service, validator, controller)
func (l *Log) SetComponent(component string) *Log {
	if component == "" {
		delete(l.Data, "component")
		l.component = ""
	} else {
		l.component = component
		l.Entry = l.WithField("component", l.component)
	}

	return l
}

// SetOperation adds the operation create/update/delete
func (l *Log) SetOperation(operation string) *Log {
	if operation == "" {
		delete(l.Data, "operation")
		l.operation = ""
	} else {
		l.operation = operation
		l.Entry = l.WithField("operation", l.operation)
	}

	return l
}

// SetObjectName adds the object name
func (l *Log) SetObjectName(objectName string) *Log {
	if objectName == "" {
		delete(l.Data, "objectName")
		l.objectName = ""
	} else {
		l.objectName = objectName
		l.Entry = l.WithField("objectName", l.objectName)
	}

	return l
}

// SetObjectState adds the state
func (l *Log) SetObjectState(state string) *Log {
	if state == "" {
		delete(l.Data, "objectState")
		l.objectState = ""
	} else {
		l.objectState = state
		l.Entry = l.WithField("objectState", l.objectState)
	}

	return l
}

// SetUser adds the user
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

// SetStep adds the step (step1, step2)
func (l *Log) SetStep(step string) *Log {
	if step == "" {
		delete(l.Data, "step")
		l.step = ""
	} else {
		l.step = step
		l.Entry = l.WithField("step", l.step)
	}

	return l
}

// SetStepState adds the phase
func (l *Log) SetStepState(state string) *Log {
	if state == "" {
		delete(l.Data, "stepState")
		l.stepState = ""
	} else {
		l.stepState = state
		l.Entry = l.WithField("stepState", l.stepState)
	}

	return l
}

// LogAuditAPI log api request and response with fields
func (l *Log) LogAuditAPI(httpType, endpoint, request, response string, responseCode int) {
	l.WithField("httpType", httpType).WithField("endpoint", endpoint).WithField(
		"request", request).WithField("responseCode", responseCode).WithField(
		"response", response).WithField("auditType", "api").Debug("audit api")

	for _, i := range []string{"auditType", "httpType", "endpoint", "request", "responseCode", "response"} {
		delete(l.Data, i)
	}
}

// LogAuditObject log object and object diffs
func (l *Log) LogAuditObject(objects ...interface{}) {
	var (
		oldObject interface{}
		newObject interface{}
		objDiff   string
		val1      = 1
	)

	if len(objects) > 0 {
		tmp, err := json.Marshal(objects[0])
		if err != nil {
			oldObject = objects[0]
		} else {
			oldObject = string(tmp)
		}
	} else {
		oldObject = "no object"
	}

	if len(objects) > val1 {
		tmp, err := json.Marshal(objects[1])
		if err != nil {
			newObject = objects[1]
		} else {
			newObject = string(tmp)
		}

		ch, _ := diff.GetDiffChangelog(objects[0], objects[1])
		if ch != nil {
			for _, c := range *ch {
				objDiff = objDiff + fmt.Sprintf("(%v, %v, %v, %v)\n", c.Path, c.Type, c.From, c.To)
			}
		}
	} else {
		newObject = "no object"
		objDiff = "no second object"
	}

	l.WithField("oldObject", oldObject).WithField("newObject", newObject).WithField(
		"objectDiff", objDiff).WithField("auditType", "object").Debug("audit object")

	for _, i := range []string{"auditType", "oldObject", "newObject", "objectDiff"} {
		delete(l.Data, i)
	}
}

// LogAuditEvent log events
func (l *Log) LogAuditEvent(message string) {
	l.WithField("auditType", "event").Debug(message)
	delete(l.Data, "auditType")
}

// SetFormatterType set format
func (l *Log) SetFormatterType(fType FormatterType) *Log {
	if fType == JSONFormatterType {
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: RFC3339NanoFixed,
		})
	} else if fType == TextFormatterType {
		l.logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			TimestampFormat:  RFC3339NanoFixed,
		})
	} else {
		l.logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			TimestampFormat:  RFC3339NanoFixed,
		})
	}

	return l
}

// SetLogFileFormatterType set file format
func (l *Log) SetLogFileFormatterType(fType FormatterType) *Log {
	if l.logFile != "" {
		rotateFileHook, err := rotatefilehook.NewRotateFileHook(getRotateConfig(l, fType))
		if err == nil {
			l.logger.AddHook(rotateFileHook)
		}
	}

	return l
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
	// Now, current logger is unchanged except for a new push/pop stack, and the current whole context (including
	// stack) is saved in the saved stack.
	l.contextStack = stack.New()
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

// PushPop function
// Run the function within a push/pop
func (l *Log) PushPop(f func()) {
	l.PushContext()
	f()
	l.PopContext()
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
func NewLoggerWithFile(filename string, maxSize, maxAge, maxBackup int) CommonLog {
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

	l := newLog(logger, stack1, stack2)
	l.logFile = filename
	l.logFileMaxSize = maxSize
	l.logFileMaxAge = maxAge
	l.logFileMaxBackup = maxBackup

	// log file
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(getRotateConfig(l, TextFormatterType))
	if err == nil {
		l.logger.AddHook(rotateFileHook)
	}

	return l
}

// NewLogger creates a logger context for a newly created logging output.
func NewLogger() CommonLog {
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

func getRotateConfig(l *Log, fType FormatterType) rotatefilehook.RotateFileConfig {
	rc := rotatefilehook.RotateFileConfig{
		Filename:   l.logFile,
		MaxSize:    l.logFileMaxSize,
		MaxBackups: l.logFileMaxBackup,
		MaxAge:     l.logFileMaxAge,
		Level:      logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: RFC3339NanoFixed,
		},
	}

	if fType == TextFormatterType {
		rc.Formatter = &logrus.TextFormatter{
			TimestampFormat: RFC3339NanoFixed,
		}
	}

	return rc
}

// Doesn't copy the stack, just the fields
func (l *Log) copyContextFrom(from *Log) *Log {
	l.logger = from.logger
	l.Entry = from.logger.WithFields(logrus.Fields{})
	l.SetCluster(from.cluster)
	l.SetApplication(from.application)
	l.SetResource(from.resource)
	l.SetComponent(from.component)
	l.SetOperation(from.operation)
	l.SetObjectName(from.objectName)
	l.SetObjectState(from.objectState)
	l.SetUser(from.user)
	l.SetStep(from.step)
	l.SetStepState(from.stepState)

	return l
}

func (l *Log) clear() {
	l.Entry = l.logger.WithFields(logrus.Fields{})

	l.cluster = ""
	l.application = ""
	l.resource = ""
	l.component = ""
	l.operation = ""
	l.objectName = ""
	l.objectState = ""
	l.user = ""
	l.step = ""
	l.stepState = ""
}
