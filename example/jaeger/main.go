package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/tongium/common-go/pkg/tracing"
)

func skipper(r *http.Request) bool {
	return strings.HasPrefix(r.RequestURI, "/healthz")
}

func main() {
	port := "8080"

	tracer, closer, err := tracing.JaegerTracer("example-service")
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	cfg := &tracing.MiddlewareConfig{
		UserIDHeaderKey:    "X-User-ID",
		RequestIDHeaderKey: "X-Request-ID",
		Skipper:            skipper,
		LogHeaderKeys:      []string{"User-Agent"},
		LogCookieKeys:      []string{"access_token"},
	}

	// Get middleware from custom configuration
	opentracingMiddleware := tracing.OpentracingMiddlewareWithConfig(cfg)

	// Simple request ID middleware
	requestIDMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerKey := "X-Request-ID"
			rid := r.Header.Get(headerKey)
			if rid == "" {
				rid = fmt.Sprintf("%d.%s", time.Now().Nanosecond(), r.RemoteAddr)
			}

			w.Header().Set(headerKey, rid)
			next.ServeHTTP(w, r)
		})
	}

	httpClient := &http.Client{}

	root := func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "root")
		defer span.Finish()

		wait(ctx)

		httpReq, _ := http.NewRequest("GET", "http://localhost:"+port+"/sleep", nil)
		tracing.SetTracingHeader(httpReq, ctx, span)
		_, err := httpClient.Do(httpReq)
		if err != nil {
			span.SetBaggageItem("error", err.Error())
			log.Println(err)
		}

		fmt.Fprintf(w, "Hello")
		w.WriteHeader(http.StatusOK)
	}

	sleep := func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContext(r.Context(), "sleep")
		defer span.Finish()

		time.Sleep(1000 * time.Millisecond)
		w.WriteHeader(http.StatusNoContent)
	}

	mux := http.NewServeMux()
	mux.Handle("/", requestIDMiddleware(opentracingMiddleware(http.HandlerFunc(root))))
	mux.Handle("/sleep", requestIDMiddleware(opentracingMiddleware(http.HandlerFunc(sleep))))

	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func wait(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "wait")
	defer span.Finish()

	time.Sleep(400 * time.Millisecond)
}
