package tracing

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"github.com/tongium/common-go/pkg/constant"
)

func ContextWithRequestID(c echo.Context) context.Context {
	ctx := c.Request().Context()
	rid := c.Request().Header.Get(echo.HeaderXRequestID)
	if rid == "" {
		rid = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	return context.WithValue(ctx, constant.RequestIDContextKey, rid)
}

func GetRequestIDFromContext(ctx context.Context) string {
	if ctx != nil && ctx.Value("request_id") != nil {
		return fmt.Sprintf("%v", ctx.Value(constant.RequestIDContextKey))
	}

	return ""
}
