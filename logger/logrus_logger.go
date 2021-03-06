package logger

import (
	"context"
	"io"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusLogger logrus.Logger implementation of interface
type LogrusLogger struct {
	Logger
	entry *logrus.Entry
	// CorrelationID is a unique identifier of the request.
	CorrelationID string
	// TenantID is a unique identifier of the tenant.
	TenantID string
}

// SetArea logger context (Area)
func (l LogrusLogger) SetArea(area string) Logger {
	return LogrusLogger{entry: l.entry.
		WithField("Area", area)}
}

// SetOperation logger context (Operation)
func (l LogrusLogger) SetOperation(operation string) Logger {
	return LogrusLogger{entry: l.entry.
		WithField("Operation", operation)}
}

// SetTenantID set TenantID for current logger context
func (l LogrusLogger) SetTenantID(tenantID string) Logger {
	if tenantID == "" {
		return l
	}
	return LogrusLogger{
		TenantID: tenantID,
		entry: l.entry.
			WithField("TenantID", tenantID),
	}
}

// SetCorrelationID set CorrelationID for logger context
func (l LogrusLogger) SetCorrelationID(correlationID string) Logger {
	if correlationID == "" {
		return l
	}
	return LogrusLogger{
		CorrelationID: correlationID,
		entry: l.entry.
			WithField("CorrelationID", correlationID),
	}
}

// GetCorrelationID returns CorrelationID for current logger context
func (l LogrusLogger) GetCorrelationID() string {
	return l.CorrelationID
}

// GetTenantID returns TenantID for current logger context
func (l LogrusLogger) GetTenantID() string {
	return l.TenantID
}

// WithField adds a filed to log entry
func (l LogrusLogger) WithField(key string, value interface{}) Logger {
	return LogrusLogger{entry: l.entry.WithField(key, value)}
}

// WithFields adds a struct of fields to the log entry.
func (l LogrusLogger) WithFields(fields Fields) Logger {
	return LogrusLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

// WithError Add an error as single field to the log entry.
func (l LogrusLogger) WithError(err error) Logger {
	return LogrusLogger{entry: l.entry.WithError(err)}
}

// WithContext adds a context to the Entry.
func (l LogrusLogger) WithContext(ctx context.Context) Logger {
	corrID := GetRequestID(ctx)
	tenantID := GetTenantID(ctx)
	log := l.SetCorrelationID(corrID).
		SetTenantID(tenantID).(LogrusLogger).
		entry.WithContext(ctx)

	return LogrusLogger{entry: log}
}

// Trace creates logs entry with Trace level
func (l LogrusLogger) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

// Debug creates logs entry with Debug level
func (l LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Info creates logs entry with Info level
func (l LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Warn creates logs entry with Warn level
func (l LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Error creates logs entry with Error level
func (l LogrusLogger) Error(args ...interface{}) {
	l.addCallerInfo().Error(args...)
}

// Fatal creates logs entry with Fatal level
func (l LogrusLogger) Fatal(args ...interface{}) {
	l.addCallerInfo().Fatal(args...)
}

// Panic creates logs entry with Panic level
func (l LogrusLogger) Panic(args ...interface{}) {
	l.addCallerInfo().Panic(args...)
}

// SetOutput sets the output to desired io.Writer like file, stdout, stderr etc
func (l *LogrusLogger) SetOutput(w io.Writer) {
	l.entry.Logger.Out = w
}

// SetLevel sets logger level
func (l LogrusLogger) SetLevel(level Level) error {
	lvl, err := logrus.ParseLevel(ParseLevel(level))
	if err != nil {
		return err
	}

	l.entry.Logger.Level = lvl
	return nil
}

// GetLevel returns current logging level
func (l LogrusLogger) GetLevel() Level {
	return ParseLogrusLevel(l.entry.Logger.Level)
}

func (l LogrusLogger) addCallerInfo() *logrus.Entry {
	// Skip this function, and fetch the PC and file for its parent.
	pc, file, line, _ := runtime.Caller(2)
	// Retrieve a function object this functions parent.

	// Regex to extract just the function name (and not the module path).
	//runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	//fname := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	funcName := runtime.FuncForPC(pc).Name()
	return l.entry.
		WithField("File", file).
		WithField("Line", line).
		WithField("Func", funcName)
}

// TrackFuncTime creates log record with func execution time
// require debug level or higher
// usage:
// func SomeFunction(list *[]string) {
//    defer TimeTrack(time.Now())
//    // Do whatever you want.
// }
func (l LogrusLogger) TrackFuncTime(start time.Time) {
	lvl := l.GetLevel()
	if lvl > DebugLevel {
		return
	}
	elapsed := time.Since(start)

	l.addCallerInfo().
		WithField("executionTime", elapsed).
		Debug("func execution completed")
}

// NewLogrusLogger creates new instance of LogrusLogger
func NewLogrusLogger(l logrus.Level) LogrusLogger {
	log := logrus.New()
	log.SetLevel(l)

	return LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// NewNopLogger returns a logger that discards all log messages.
func NewNopLogger() Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	return LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// ParseLogrusLevel takes a string level and returns the Logrus log level constant.
func ParseLogrusLevel(lvl logrus.Level) Level {
	switch lvl {
	case logrus.PanicLevel:
		return PanicLevel
	case logrus.FatalLevel:
		return FatalLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.TraceLevel:
		return TraceLevel
	}
	// Default value
	return InfoLevel
}
