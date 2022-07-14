package log

type noopLogger struct{}

func (l *noopLogger) Info(_ string, _ ...interface{})  {}
func (l *noopLogger) Error(_ string, _ ...interface{}) {}
func (l *noopLogger) Debug(_ string, _ ...interface{}) {}

func NewNoop() Interface {
	return &noopLogger{}
}
