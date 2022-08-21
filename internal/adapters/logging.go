package adapters

type Logger interface {
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
	Debugf(format string, params ...interface{})
	Tracef(format string, params ...interface{})
	Panicf(format string, params ...interface{})
}

// AppLogger логгер всего приложения
type AppLogger struct {
	l Logger
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
	a.l.Errorf(format, params)
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
