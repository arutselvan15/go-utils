package log

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const (
	RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"
	DebugLevel       = "debug"
	InfoLevel        = "info"
)

type logFormat string

var FormatText logFormat = "Text"
var FormatJson logFormat = "Json"

type log struct {
	*logrus.Entry
	entry    *logrus.Entry
	action   string
	step     string
	endpoint string
	payload  string
	response string
}

type GoLog interface {
	logrus.FieldLogger
	SetLevel(string) *log
	SetAction(action string) *log
	SetStep(step string) *log
	SetAPI(string, string, string) *log
	GetLogger() *log
}

// SetLevel sets the level at which log messages are published/written.
func (l *log) SetLevel(level string) *log {
	// If there's no explicit logging level specified, set the level to INFO
	if level == "" {
		level = InfoLevel
	}

	loglevel, err := logrus.ParseLevel(level)
	if err == nil {
		// set default logger and the custom logger levels
		logrus.SetLevel(loglevel)
		l.entry.Logger.SetLevel(loglevel)
	}
	return l
}

// SetAction adds the action (Create, Update, Delete, etc.) field to each log message if provided
func (l *log) SetAction(action string) *log {
	if action != "" {
		l.action = action
		l.Entry = l.WithField("action", l.action)
	}
	return l
}

// SetStep adds the step into the log (step1, step2 and etc.)
func (l *log) SetStep(step string) *log {
	if step != "" {
		l.step = step
		l.Entry = l.WithField("step", l.step)
	}
	return l
}

// SetAPI sets the endpoint, payload, and response of an api call
func (l *log) SetAPI(endpoint, payload, response string) *log {
	if endpoint != "" {
		l.endpoint = endpoint
		l.Entry = l.WithField("endpoint", l.endpoint)
	}

	if payload != "" {
		l.payload = payload
		l.Entry = l.WithField("payload", l.payload)
	}

	if response != "" {
		l.response = response
		l.Entry = l.WithField("response", l.response)
	}
	return l
}

// GetLogger returns the logrus object
func (l *log) GetLogger() *log {
	l.Entry = l.entry
	l.action = ""
	l.step = ""
	l.endpoint = ""
	l.payload = ""
	l.response = ""
	return l
}

// NewLogger is the constructor for Log
func NewLogger(component string, format logFormat, logFilePath string) *log {
	l := &log{}
	logger := logrus.New()

	var formatter logrus.Formatter

	if format == FormatJson {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: RFC3339NanoFixed,
		}
	} else {
		// default text formatter
		formatter = &logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			TimestampFormat:  RFC3339NanoFixed,
		}
	}

	// log file update config
	if logFilePath != "" {
		rotateConfig := rotatefilehook.RotateFileConfig{
			Filename:   logFilePath,
			MaxSize:    250,
			MaxBackups: 3,
			MaxAge:     30,
			Level:      logrus.DebugLevel,
			Formatter:  formatter,
		}

		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotateConfig)
		if err != nil {
			logrus.Fatalf("failed to initialize file rotate hook: %v", err)
		}
		logger.AddHook(rotateFileHook)
	}

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(formatter)

	l.entry = logger.WithFields(logrus.Fields{
		"component": component,
	})
	l.Entry = l.entry

	return l
}
