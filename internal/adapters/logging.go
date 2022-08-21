package adapters

type Logger interface {
	Infof(format string, params ...interface{})
	Info(args ...interface{})
	Warnf(format string, params ...interface{})
	Warn(args ...interface{})
	Errorf(format string, params ...interface{})
	Error(args ...interface{})
	Debugf(format string, params ...interface{})
	Debug(args ...interface{})
	Tracef(format string, params ...interface{})
	Trace(args ...interface{})
	Panicf(format string, params ...interface{})
	Panic(args ...interface{})
}

// AppLogger логгер всего приложения
type AppLogger struct {
	l Logger
}

func (a *AppLogger) Info(args ...interface{}) {
	a.l.Info(args)
}

func (a *AppLogger) Warn(args ...interface{}) {
	a.l.Warn(args)
}

func (a *AppLogger) Error(args ...interface{}) {
	a.l.Error(args)
}

func (a *AppLogger) Debug(args ...interface{}) {
	a.l.Debug(args)
}

func (a *AppLogger) Trace(args ...interface{}) {
	a.l.Trace(args)
}

func (a *AppLogger) Panic(args ...interface{}) {
	a.l.Panic(args)
}

func GetAppLogger(l *Logger) *AppLogger {
	return &AppLogger{*l}
}

func (a *AppLogger) Infof(format string, params ...interface{}) {
	a.l.Infof(format, params)
}

func (a *AppLogger) Warnf(format string, params ...interface{}) {
	a.l.Warnf(format, params)
}

func (a *AppLogger) Errorf(format string, params ...interface{}) {
	a.l.Errorf(format, params...)
}

func (a *AppLogger) Debugf(format string, params ...interface{}) {
	a.l.Debugf(format, params)
}

func (a *AppLogger) Tracef(format string, params ...interface{}) {
	a.l.Tracef(format, params)
}

func (a *AppLogger) Panicf(format string, params ...interface{}) {
	a.l.Panicf(format, params)
}
