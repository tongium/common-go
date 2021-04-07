package tracing

import (
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

func GetSpanByRequest(functionName string, r *http.Request) opentracing.Span {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	return tracer.StartSpan(functionName, ext.RPCServerOption(spanCtx))
}

func SetSpanRequest(span opentracing.Span, req *http.Request) error {
	tracer := opentracing.GlobalTracer()
	return tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}
