package tracing

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/tongium/common-go/pkg/constant"
)

type responseWriter struct {
	http.ResponseWriter
	span opentracing.Span
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.span.SetTag("http.status_code", code)
	rw.ResponseWriter.WriteHeader(code)
}

type MiddlewareConfig struct {
	Tracer             opentracing.Tracer
	RequestIDHeaderKey string
	UserIDHeaderKey    string
	Skipper            func(r *http.Request) bool
	LogHeaderKeys      []string
	LogCookieKeys      []string
}

func DefaultSkipper(r *http.Request) bool {
	return false
}

func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		Tracer:             opentracing.GlobalTracer(),
		RequestIDHeaderKey: constant.DefaultRequstIDHeaderKey,
		UserIDHeaderKey:    "",
		Skipper:            DefaultSkipper,
	}
}

// Get middleware with config
func OpentracingMiddlewareWithConfig(cfg *MiddlewareConfig) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultMiddlewareConfig()
	}

	if cfg.Tracer == nil {
		cfg.Tracer = opentracing.GlobalTracer()
	}

	if cfg.Skipper == nil {
		cfg.Skipper = DefaultSkipper
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			span := getSpanByRequest(cfg.Tracer, r)
			defer span.Finish()

			span.SetTag("http.method", r.Method)
			span.SetTag("http.url", r.URL)
			span.SetTag("http.scheme", r.Proto)

			if len(cfg.LogHeaderKeys) > 0 {
				addHeaderLog(cfg.LogHeaderKeys, r, span)
			}

			if len(cfg.LogCookieKeys) > 0 {
				addCookieLog(cfg.LogCookieKeys, r, span)
			}

			writer := &responseWriter{
				ResponseWriter: w,
				span:           span,
			}

			ctx := r.Context()
			if cfg.RequestIDHeaderKey != "" {
				if rid := getRequestID(cfg.RequestIDHeaderKey, r, writer); rid != "" {
					span.SetTag("http.request_id", rid)
					ctx = context.WithValue(r.Context(), constant.RequestIDContextKey, rid)
				}
			}

			if cfg.UserIDHeaderKey != "" {
				if userID := r.Header.Get(cfg.UserIDHeaderKey); userID != "" {
					span.SetTag("http.user_id", userID)
				}
			}

			req := r.WithContext(opentracing.ContextWithSpan(ctx, span))
			next.ServeHTTP(writer, req)
		})
	}
}

// Get middleware with default config
func OpentracingMiddleware() func(http.Handler) http.Handler {
	return OpentracingMiddlewareWithConfig(DefaultMiddlewareConfig())
}

func addHeaderLog(keys []string, r *http.Request, span opentracing.Span) {
	if r == nil || span == nil {
		return
	}

	for _, key := range keys {
		if value := r.Header.Get(key); value != "" {
			span.LogFields(log.String("header:"+strings.ToLower(key), r.Header.Get(key)))
		}
	}
}

func addCookieLog(keys []string, r *http.Request, span opentracing.Span) {
	if r == nil || span == nil {
		return
	}

	for _, key := range keys {
		if c, err := r.Cookie(key); err == nil {
			if value := c.Value; value != "" {
				span.LogFields(log.String("cookie:"+strings.ToLower(key), value))
			}
		}
	}
}

func getSpanByRequest(tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	operationName := fmt.Sprintf("%s %s", r.Method, r.RequestURI)

	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		return tracer.StartSpan(operationName)
	}

	return tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
}

func getRequestID(key string, r *http.Request, w http.ResponseWriter) string {
	rid := r.Header.Get(key)
	if rid == "" {
		rid = w.Header().Get(key)
	}

	return rid
}
