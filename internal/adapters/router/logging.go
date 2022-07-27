package router

type Logger interface {
	Infof(format string, params ...interface{})
}
