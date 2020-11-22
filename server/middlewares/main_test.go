package middlewares

type loggerCallArgs struct {
	f string
	v []interface{}
}

type mockLogger struct {
	info  loggerCallArgs
	debug loggerCallArgs
	warn  loggerCallArgs
	error loggerCallArgs
}

func (l *mockLogger) Debug(format string, v ...interface{}) {
	l.debug.f = format
	l.debug.v = v
}

func (l *mockLogger) Info(format string, v ...interface{}) {
	l.info.f = format
	l.info.v = v
}

func (l *mockLogger) Warn(format string, v ...interface{}) {
	l.warn.f = format
	l.warn.v = v
}

func (l *mockLogger) Error(format string, v ...interface{}) {
	l.error.f = format
	l.error.v = v
}
