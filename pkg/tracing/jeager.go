package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/tongium/common-go/pkg/constant"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type Config struct {
	Tracer             opentracing.Tracer
	RequestIDHeaderKey string
}

func DefaultTracerConfig() *Config {
	return &Config{
		Tracer:             opentracing.GlobalTracer(),
		RequestIDHeaderKey: constant.DefaultRequstIDHeaderKey,
	}
}

// Get tracer with default configuration if environment not found.
// https://github.com/jaegertracing/jaeger-client-go#environment-variables
func JaegerTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	defaultConfig := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	// Override default configuration
	config, err := defaultConfig.FromEnv()
	if err != nil {
		return nil, nil, err
	}

	return config.NewTracer()
}

// Carries span and request-id to request header
func SetTracingHeader(req *http.Request, ctx context.Context, span opentracing.Span) {
	if req == nil || ctx == nil || span == nil {
		return
	}

	SetTracingHeaderWithConfig(req, ctx, span, DefaultTracerConfig())
}

// Carries span and request-id to request header with custom configuration
func SetTracingHeaderWithConfig(req *http.Request, ctx context.Context, span opentracing.Span, cfg *Config) {
	if cfg == nil {
		cfg = DefaultTracerConfig()
	}

	if cfg.Tracer == nil {
		cfg.Tracer = opentracing.GlobalTracer()
	}

	if ctx != nil && cfg.RequestIDHeaderKey != "" {
		if rid := GetRequestIDFromContext(ctx); rid != "" {
			req.Header.Set(cfg.RequestIDHeaderKey, rid)
		}
	}

	if span != nil {
		cfg.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	}
}

// Get request id from context which set by middleware
func GetRequestIDFromContext(ctx context.Context) string {
	if ctx != nil && ctx.Value(constant.RequestIDContextKey) != nil {
		return fmt.Sprintf("%v", ctx.Value(constant.RequestIDContextKey))
	}

	return ""
}
