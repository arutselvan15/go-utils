package log

import (
	"fmt"
	"os"

	"github.com/golang-collections/collections/stack"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"gopkg.in/yaml.v2"

	"github.com/arutselvan15/go-utils/diff"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

type log struct {
	*logrus.Entry

	// creator
	logger *logrus.Logger

	// used for push/pop of contexts
	contextStack *stack.Stack

	// save/restore context stack
	savedContexts *stack.Stack

	// fields
	component       string
	subcomponent    string
	process         string
	subprocess      string
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

type UtilsLog interface {
	logrus.FieldLogger
	SetLevel(string) *log
	SetSubComponent(string) *log
	SetProcess(process string) *log
	SetSubProcess(subProcess string) *log
	SetAction(action string) *log
	SetUser(user string) *log
	SetInvolvedObj(involvedObj string) *log
	SetDisposition(disposition string) *log
	SetAPIRequest(string, string) *log
	SetAPIResponse(string, string) *log
	Event(string, ...interface{})
	SetObjectAudit(...interface{}) *log
	GetLogger() *log
	PushContext()
	PopContext()
	SaveContext()
	RestoreContext()
	PushPop(f func())
	ThreadLogger() *log
}

// SetObjectAudit sets the object audit data and type
func (l *log) SetObjectAudit(objects ...interface{}) *log {
	if len(objects) == 1 {
		// create
		l.objectAuditType = "create"
		l.Entry = l.WithField("object_audit_type", l.objectAuditType)
		yamlBytes, _ := yaml.Marshal(objects[0])
		l.objectAuditData = string(yamlBytes)
		l.Entry = l.WithField("object_audit_data", l.objectAuditData)

	} else if len(objects) == 2 {
		// update and look at first 2
		l.objectAuditType = "update"
		l.Entry = l.WithField("object_audit_type", l.objectAuditType)
		ch, _ := diff.GetDiffChangelog(objects[0], objects[1])
		data := ""
		if ch != nil {
			for _, c := range *ch {
				data = data + fmt.Sprintf("(%v, %v, %v, %v)\n", c.Path, c.Type, c.From, c.To)
			}
			l.objectAuditData = data
			l.Entry = l.WithField("object_audit_data", l.objectAuditData)
		}
	}

	return l
}

// SetLevel sets the level at which log messages are published/written.
func (l *log) SetLevel(level string) *log {
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

// SetAPI sets the endpoint, payload, and response of an api call
func (l *log) SetAPIRequest(endpoint, request string) *log {

	if endpoint == "" {
		delete(l.Data, "api_endpoint")
		l.apiEndpoint = ""
	} else {
		l.apiEndpoint = endpoint
		l.Entry = l.WithField("api_endpoint", l.apiEndpoint)
	}

	if request == "" {
		delete(l.Data, "api_request")
		l.apiRequest = ""
	} else {
		l.apiRequest = request
		l.Entry = l.WithField("api_request", l.apiRequest)
	}

	return l
}

func (l *log) SetAPIResponse(endpoint, response string) *log {

	if endpoint == "" {
		delete(l.Data, "api_endpoint")
		l.apiEndpoint = ""
	} else {
		l.apiEndpoint = endpoint
		l.Entry = l.WithField("api_endpoint", l.apiEndpoint)
	}

	if response == "" {
		delete(l.Data, "api_response")
		l.apiResponse = ""
	} else {
		l.apiResponse = response
		l.Entry = l.WithField("api_response", l.apiResponse)
	}

	return l
}

// SetSubComponent adds the sub component (ping, am, cert, etc.) field to to each log message if provided
func (l *log) SetSubComponent(sc string) *log {
	if sc == "" {
		delete(l.Data, "subcomponent")
		l.subcomponent = ""
	} else {
		l.subcomponent = sc
		l.Entry = l.WithField("subcomponent", l.subcomponent)
	}

	return l
}

// SetProcess adds the process (CreateProject, DeleteProject, etc.) field to each log message if provided
func (l *log) SetProcess(process string) *log {
	if process == "" {
		delete(l.Data, "process")
		l.process = ""
	} else {
		l.process = process
		l.Entry = l.WithField("process", l.process)
	}

	return l
}

// SetSubProcess adds the subprocess field to each log message if provided
func (l *log) SetSubProcess(subProcess string) *log {
	if subProcess == "" {
		delete(l.Data, "subprocess")
		l.subprocess = ""
	} else {
		l.subprocess = subProcess
		l.Entry = l.WithField("subprocess", l.subprocess)
	}

	return l
}

// SetAction adds the action into the log (configureProjectAdmin, configureRBAC and etc.)
func (l *log) SetAction(action string) *log {
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
func (l *log) SetUser(user string) *log {
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
func (l *log) SetInvolvedObj(involvedObj string) *log {
	if involvedObj == "" {
		delete(l.Data, "involved_object")
		l.involvedObj = ""
	} else {
		l.involvedObj = involvedObj
		l.Entry = l.WithField("involved_object", l.involvedObj)
	}

	return l
}

// SetDisposition adds the disposition (Success/Fail of the process and/or action) field to each log message if provided
func (l *log) SetDisposition(disposition string) *log {
	if disposition == "" {
		delete(l.Data, "disposition")
		l.disposition = ""
	} else {
		l.disposition = disposition
		l.Entry = l.WithField("disposition", l.disposition)
	}

	return l
}

func (l *log) Event(format string, args ...interface{}) {
	l.PushContext()
	l.SetAction("Event").Infof(format, args)
	l.PopContext()
}

// Save the current context for later restore.
// This saves the current fields and the current stack.
// After this call, the current context is intact, but has a new empty stack.
func (l *log) SaveContext() {
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

// If previously saved, this restores the saved context.
// This restores the previouosly saved field and the previous stack.
// Anything in the stack between Save/Restore is gone.
func (l *log) RestoreContext() {
	// If no saves, do nothing
	if l.savedContexts.Len() == 0 {
		return
	}
	// Get the last saved full context
	c := l.savedContexts.Pop().(*log)
	// Restore saved push/pop stack
	l.contextStack = c.contextStack
	// Restore context data
	l.copyContextFrom(c)
}

func (l *log) PushContext() {
	// push and pop by value, not by reference
	l.contextStack.Push(*l)
}

func (l *log) PopContext() {
	// Do nothing if nothing there
	if l.contextStack.Len() == 0 {
		return
	}

	pop := l.contextStack.Pop().(log)

	l.copyContextFrom(&pop)
}

// Run the function within a push/pop
func (l *log) PushPop(f func()) {
	l.PushContext()
	f()
	l.PopContext()
}

// Doesn't copy the stack, just the fields
func (l *log) copyContextFrom(from *log) *log {
	l.component = from.component
	l.logger = from.logger
	l.Entry = from.logger.WithFields(logrus.Fields{
		"cluster":   os.Getenv("CLUSTER"),
		"component": from.component,
	})
	l.SetSubComponent(from.subcomponent)
	l.SetProcess(from.process)
	l.SetSubProcess(from.subprocess)
	l.SetAction(from.action)
	l.SetUser(from.user)
	l.SetInvolvedObj(from.involvedObj)
	l.SetDisposition(from.disposition)
	l.SetAPIRequest(from.apiEndpoint, from.apiRequest)
	l.SetAPIResponse(from.apiEndpoint, from.apiResponse)
	l.setObjectAudit(l.objectAuditType, l.objectAuditData)

	return l
}

func (l *log) setObjectAudit(auditType string, data string) {
	if auditType != "" {
		l.objectAuditType = auditType
		l.Entry = l.WithField("object_audit_type", l.objectAuditType)
	}

	if data != "" {
		l.objectAuditData = data
		l.Entry = l.WithField("object_audit_data", l.objectAuditData)
	}
}

func (l *log) clear() {
	l.Entry = l.logger.WithFields(logrus.Fields{
		"cluster":   os.Getenv("CLUSTER"),
		"component": l.component,
	})

	l.subcomponent = ""
	l.process = ""
	l.subprocess = ""
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

func newlog(logger *logrus.Logger, component string, contextStack *stack.Stack, savedContexts *stack.Stack) *log {
	nl := &log{
		component:     component,
		logger:        logger,
		contextStack:  contextStack,
		savedContexts: savedContexts,
	}

	nl.clear()
	return nl
}

// GetLogger creates and returns a new logging context.
// A logging context is a wrapper for a log entry for a logger.  These objects can NOT be used
// in parallel.
func (l *log) GetLogger() *log {
	nl := newlog(l.logger, l.component, l.contextStack, l.savedContexts)
	return nl
}

// ThreadLogger clones the logger, but gives it a new set of stacks.
// This allows the thread to do all context stuff independenly of other threads
func (l *log) ThreadLogger() *log {
	// create the context stacks
	stack1 := stack.New()
	stack2 := stack.New()

	// create the initial logging context
	nlog := newlog(l.logger, l.component, stack1, stack2)
	// populate with existing context
	nlog.copyContextFrom(l)
	return nlog
}

func NewLoggerWithFile(component string, filename string) *log {
	logger := logrus.New()

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

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  RFC3339NanoFixed,
	})
	logger.AddHook(rotateFileHook)

	// create the shared context stacks
	stack1 := stack.New()
	stack2 := stack.New()

	// create the initial logging context
	nlog := newlog(logger, component, stack1, stack2)

	return nlog
}

// NewLogger creates a logger context for a newly created logging output.
func NewLogger(component string) *log {
	return NewLoggerWithFile(component, "/log/felix-"+component+".log")
}
