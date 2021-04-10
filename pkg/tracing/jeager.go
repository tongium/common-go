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

func getDefaultTracer() *Config {
	return &Config{
		Tracer:             opentracing.GlobalTracer(),
		RequestIDHeaderKey: constant.DefaultRequstIDHeaderKey,
	}
}

func GetJaegerTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	return cfg.NewTracer()
}

func inject(tracer opentracing.Tracer, span opentracing.Span, req *http.Request) error {
	return tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}

// Carries span and request-id to request header
func SetTracingHeader(ctx context.Context, req *http.Request, span opentracing.Span, cfg *Config) {
	if cfg == nil {
		cfg = getDefaultTracer()
	}

	if cfg.Tracer == nil {
		cfg.Tracer = opentracing.GlobalTracer()
	}

	if ctx != nil && cfg.RequestIDHeaderKey != "" {
		if rid := GetRequestIDFromContext(ctx); rid != "" {
			req.Header.Set(cfg.RequestIDHeaderKey, rid)
		}
	}

	if cfg.Tracer != nil && span != nil {
		inject(cfg.Tracer, span, req)
	}
}

func GetRequestIDFromContext(ctx context.Context) string {
	if ctx != nil && ctx.Value(constant.RequestIDContextKey) != nil {
		return fmt.Sprintf("%v", ctx.Value(constant.RequestIDContextKey))
	}

	return ""
}
