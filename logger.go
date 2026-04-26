package cache

type Logger interface {
	Errorf(string, ...interface{})
}

// NoopLogger the default logger that will discard all logs of gin-cache
type NoopLogger struct {
}

// Errorf will output the log at error level
func (l NoopLogger) Errorf(string, ...interface{}) {
}
