package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tongium/common-go/pkg/constant"
)

type responseWriter struct {
	http.ResponseWriter
	span        opentracing.Span
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.ResponseWriter.WriteHeader(code)
	rw.span.SetTag("http.status_code", code)
	rw.wroteHeader = true
}

type MiddlewareConfig struct {
	Tracer             opentracing.Tracer
	RequestIDHeaderKey string
	UserIDHeaderKey    string
	Skipper            func(r *http.Request) bool
}

func DefaultSkipper(r *http.Request) bool {
	return false
}

func getDefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		Tracer:             opentracing.GlobalTracer(),
		RequestIDHeaderKey: constant.DefaultRequstIDHeaderKey,
		UserIDHeaderKey:    "",
		Skipper:            DefaultSkipper,
	}
}

func JaegerMiddleware(cfg *MiddlewareConfig) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = getDefaultMiddlewareConfig()
	}

	if cfg.Tracer == nil {
		cfg.Tracer = opentracing.GlobalTracer()
	}

	if cfg.Skipper == nil {
		cfg.Skipper = DefaultSkipper
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if cfg.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			operationName := fmt.Sprintf("%s %s", r.Method, r.RequestURI)
			spanCtx, err := cfg.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			var span opentracing.Span
			if err != nil {
				span = cfg.Tracer.StartSpan(operationName)
			} else {
				span = cfg.Tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
			}

			defer span.Finish()

			writer := &responseWriter{
				ResponseWriter: w,
				span:           span,
			}

			var rid string
			ctx := r.Context()

			if cfg.RequestIDHeaderKey != "" {
				rid = r.Header.Get(cfg.RequestIDHeaderKey)
				if rid == "" {
					rid = writer.Header().Get(cfg.RequestIDHeaderKey)
				}

				if rid != "" {
					ctx = context.WithValue(r.Context(), constant.RequestIDContextKey, rid)
					span.SetTag("http.request_id", rid)
				}
			}

			if cfg.UserIDHeaderKey != "" {
				if userID := r.Header.Get(cfg.UserIDHeaderKey); userID != "" {
					span.SetTag("http.user_id", userID)
				}
			}

			req := r.WithContext(opentracing.ContextWithSpan(ctx, span))
			next.ServeHTTP(writer, req)
		}

		return http.HandlerFunc(fn)
	}
}
