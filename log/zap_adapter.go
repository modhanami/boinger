package log

import "go.uber.org/zap"

type zapLoggerAdapter struct {
	logger *zap.SugaredLogger
}

func NewZapLoggerAdapter(l *zap.SugaredLogger) *zapLoggerAdapter {
	return &zapLoggerAdapter{logger: l}
}

func (l *zapLoggerAdapter) Info(msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *zapLoggerAdapter) Error(msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *zapLoggerAdapter) Debug(msg string, args ...interface{}) {
	l.logger.Debugw(msg, args...)
}
