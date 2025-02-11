package worker

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type logCapture struct {
	Level zerolog.Level
	Msg   string
}

func captureOutput(*testing.T) (*[]logCapture, func()) {
	var output []logCapture
	originalLogger := log.Logger
	log.Logger = log.Logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		output = append(output, logCapture{Level: level, Msg: msg})
	}))

	return &output, func() {
		log.Logger = originalLogger
	}
}

func TestLogger(t *testing.T) {
	logger := NewLogger()

	tests := []struct {
		name     string
		logFunc  func()
		level    zerolog.Level
		expected string
	}{
		{
			name:     "Debug level",
			logFunc:  func() { logger.Debug("debug message") },
			level:    zerolog.DebugLevel,
			expected: "debug message",
		},
		{
			name:     "Info level",
			logFunc:  func() { logger.Info("info message") },
			level:    zerolog.InfoLevel,
			expected: "info message",
		},
		{
			name:     "Warn level",
			logFunc:  func() { logger.Warn("warn message") },
			level:    zerolog.WarnLevel,
			expected: "warn message",
		},
		{
			name:     "Error level",
			logFunc:  func() { logger.Error("error message") },
			level:    zerolog.ErrorLevel,
			expected: "error message",
		},
		{
			name:     "Fatal level",
			logFunc:  func() { logger.Fatal("fatal message") },
			level:    zerolog.FatalLevel,
			expected: "fatal message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, cleanup := captureOutput(t)
			defer cleanup()

			tt.logFunc()

			assert.Len(t, *output, 1)
			assert.Equal(t, tt.level, (*output)[0].Level)
			assert.Equal(t, tt.expected, (*output)[0].Msg)
		})
	}
}

func TestLoggerPrintf(t *testing.T) {
	logger := NewLogger()
	output, cleanup := captureOutput(t)
	defer cleanup()

	ctx := context.Background()
	logger.Printf(ctx, "test %s %d", "message", 42)

	assert.Len(t, *output, 1)
	assert.Equal(t, zerolog.DebugLevel, (*output)[0].Level)
	assert.Equal(t, "test message 42", (*output)[0].Msg)
}
