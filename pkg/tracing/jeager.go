package tracing

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tongium/common-go/pkg/constant"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func GetJeagerTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
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

func GetSpanByRequest(operationName string, r *http.Request) (opentracing.Span, context.Context) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
	ctx := opentracing.ContextWithSpan(r.Context(), span)

	rid := r.Header.Get(constant.DefaultRequstIDHeaderKey)
	if rid != "" {
		span.SetTag("http.request_id", rid)
	}

	return span, ctx
}

func GetContextAndSpanByEchoContext(c echo.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := GetSpanByRequest(operationName, c.Request())

	rid := c.Request().Header.Get(echo.HeaderXRequestID)
	if rid == "" {
		rid = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	if rid != "" {
		span.SetTag("http.request_id", rid)
	}

	return span, context.WithValue(ctx, constant.RequestIDContextKey, rid)
}

func SetSpanRequest(span opentracing.Span, req *http.Request) error {
	tracer := opentracing.GlobalTracer()
	return tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}
