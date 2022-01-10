package tracing

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
)

// Option - define options for NewJWTCache()
type Option func(*options)
type options struct {
	sampleProbability float64
	enableInfoLog     bool
}

// defaultOptions - some defs options to NewJWTCache()
var defaultOptions = options{
	sampleProbability: 0.0,
	enableInfoLog:     false,
}

// WithSampleProbability - optional sample probability
func WithSampleProbability(sampleProbability float64) Option {
	return func(o *options) {
		o.sampleProbability = sampleProbability
	}
}

// WithEnableInfoLog - optional: enable Info logging for tracing
func WithEnableInfoLog(enable bool) Option {
	return func(o *options) {
		o.enableInfoLog = enable
	}
}

// InitTracing - init opentracing with options (WithSampleProbability, WithEnableInfoLog) defaults: constant sampling, no info logging
func InitTracing(serviceName string, tracingAgentHostPort string, opt ...Option) (
	tracer opentracing.Tracer,
	reporter jaeger.Reporter,
	closer io.Closer,
	err error) {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	factory := jaegerprom.New()
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})
	transport, err := jaeger.NewUDPTransport(tracingAgentHostPort, 0)
	if err != nil {
		return tracer, reporter, closer, err
	}

	logAdapt := Log{InfoLevel: opts.enableInfoLog}
	reporter = jaeger.NewCompositeReporter(
		jaeger.NewLoggingReporter(logAdapt),
		jaeger.NewRemoteReporter(transport,
			jaeger.ReporterOptions.Metrics(metrics),
			jaeger.ReporterOptions.Logger(logAdapt),
		),
	)

	var sampler jaeger.Sampler
	sampler = jaeger.NewConstSampler(true)
	if opts.sampleProbability > 0 {
		fmt.Println("probable")
		sampler, err = jaeger.NewProbabilisticSampler(opts.sampleProbability)
	}

	tracer, closer = jaeger.NewTracer(serviceName,
		sampler,
		reporter,
		jaeger.TracerOptions.Metrics(metrics),
	)
	return tracer, reporter, closer, nil
}
func GetTraceIdFromSpanContext(spanCtx opentracing.SpanContext) string {
	sc, ok := spanCtx.(jaeger.SpanContext)
	if ok {
		return sc.TraceID().String()
	}

	return ""
}

func GetSpanFromGinContex(funcName string, c *gin.Context) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if cspan, ok := c.Get("tracing-context"); ok {
		span = StartSpanWithParent(cspan.(opentracing.Span).Context(), funcName, c.Request.Method, c.Request.URL.Path)
	} else {
		span = StartSpanWithHeader(&c.Request.Header, funcName, c.Request.Method, c.Request.URL.Path)
	}
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	return span, ctx
}
