package log

type Interface interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Debug(string, ...interface{})
}
