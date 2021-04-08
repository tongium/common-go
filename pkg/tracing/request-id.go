package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
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
	if ctx != nil && ctx.Value(constant.RequestIDContextKey) != nil {
		return fmt.Sprintf("%v", ctx.Value(constant.RequestIDContextKey))
	}

	return ""
}

func SetRequestHeader(ctx context.Context, req *http.Request, span opentracing.Span, headerKey string) {
	key := headerKey
	if key == "" {
		key = constant.DefaultRequstIDHeaderKey
	}

	if ctx != nil {
		req.Header.Set(key, GetRequestIDFromContext(ctx))
	}

	if span != nil {
		SetSpanRequest(span, req)
	}
}
