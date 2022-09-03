package logger

type noopLogger struct{}

func (l *noopLogger) Debug(_ string, _ ...interface{}) {}
func (l *noopLogger) Info(_ string, _ ...interface{})  {}
func (l *noopLogger) Warn(_ string, _ ...interface{})  {}
func (l *noopLogger) Error(_ string, _ ...interface{}) {}
func (l *noopLogger) Fatal(_ string, _ ...interface{}) {}

func NewNoopLogger() Logger {
	return &noopLogger{}
}
