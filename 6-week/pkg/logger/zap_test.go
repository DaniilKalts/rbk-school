package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/6-week/pkg/logger"
)

func TestNew_SupportedLevels(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error", "DEBUG", " info "} {
		t.Run(level, func(t *testing.T) {
			l, err := logger.New(level, "json")
			require.NoError(t, err)
			require.NotNil(t, l)
		})
	}
}

func TestNew_UnsupportedLevel(t *testing.T) {
	l, err := logger.New("trace", "json")
	require.Error(t, err)
	assert.Nil(t, l)
	assert.Contains(t, err.Error(), "неподдерживаемый уровень")
}

func TestNew_ConsoleFormat(t *testing.T) {
	l, err := logger.New("info", "console")
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := logger.New("info", "xml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "сборка логгера")
}

func TestFromContext_ReturnsNopWhenAbsent(t *testing.T) {
	got := logger.FromContext(context.Background())

	require.NotNil(t, got)
	got.Info("safe to call")
}

func TestWithContext_Roundtrip(t *testing.T) {
	base := zap.NewNop()
	ctx := logger.WithContext(context.Background(), base)

	assert.Same(t, base, logger.FromContext(ctx))
}

func TestFromContext_IgnoresNilLogger(t *testing.T) {
	ctx := logger.WithContext(context.Background(), nil)

	got := logger.FromContext(ctx)
	require.NotNil(t, got)
}
