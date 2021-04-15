package logutil

import (
	"context"
	"fmt"

	"github.com/tongium/common-go/pkg/constant"
	"github.com/tongium/common-go/pkg/tracing"
	"go.uber.org/zap"
)

// Add request ID from middleware to logging
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	if rid := tracing.GetRequestIDFromContext(ctx); rid != "" {
		return logger.With(zap.String(fmt.Sprint(constant.RequestIDContextKey), rid))
	}

	return logger
}
