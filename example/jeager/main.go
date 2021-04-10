package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/tongium/common-go/pkg/tracing"
)

func skipper(r *http.Request) bool {
	return strings.HasPrefix(r.RequestURI, "/healthz")
}

func main() {
	tracer, closer, err := tracing.GetJaegerTracer("example-service")
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	cfg := &tracing.MiddlewareConfig{
		UserIDHeaderKey:    "X-User-ID",
		RequestIDHeaderKey: "X-Request-ID",
		Skipper:            skipper,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID(), echo.WrapMiddleware(tracing.JaegerMiddleware(cfg)))

	httpClient := &http.Client{}

	e.GET("/", func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "handler")
		defer span.Finish()

		wait(ctx)

		httpReq, _ := http.NewRequest("GET", "http://localhost:1323/sleep", nil)
		tracing.SetTracingHeader(ctx, httpReq, span, nil)
		_, err := httpClient.Do(httpReq)
		if err != nil {
			span.SetBaggageItem("error", err.Error())
			c.Logger().Error(err)
		}

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/sleep", func(c echo.Context) error {
		span, _ := opentracing.StartSpanFromContext(c.Request().Context(), "handler")
		defer span.Finish()

		time.Sleep(1000 * time.Millisecond)

		return c.String(http.StatusNoContent, "")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func wait(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "wait")
	defer span.Finish()

	time.Sleep(400 * time.Millisecond)
}
