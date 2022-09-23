package logger

type noopLogger struct{}

var _ Logger = (*noopLogger)(nil)

func (l *noopLogger) Info(_ string, _ ...interface{})  {}
func (l *noopLogger) Error(_ string, _ ...interface{}) {}
func (l *noopLogger) Debug(_ string, _ ...interface{}) {}
func (l *noopLogger) Warn(_ string, _ ...interface{})  {}
func (l *noopLogger) Fatal(_ string, _ ...interface{}) {}
func (l *noopLogger) With(_ ...interface{}) Logger {
	return &noopLogger{}
}

func NewNoopLogger() Logger {
	return &noopLogger{}
}
