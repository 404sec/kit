// Utilities for working with tracing.
// dependencies:
//   go get github.com/opentracing/opentracing-go
//   go get github.com/uber/jaeger-client-go

package tracing

import (
	"context"
	"github.com/404sec/log"
	opentracing "github.com/opentracing/opentracing-go"
)

// Log - an adapter to log span info
type Log struct {
	InfoLevel bool
}

// Error - logrus adapter for span errors
func (l Log) Error(msg string) {
	log.Error(context.Background(), "Reporting Span Err", msg)
}

// Infof - logrus adapter for span info logging
func (l Log) Infof(msg string, args ...interface{}) {
	if l.InfoLevel {
		ctx := context.Background()
		for _, v := range args {
			ctx = opentracing.ContextWithSpan(ctx, v.(opentracing.Span))
		}

		log.Infof(ctx, msg, args)

	}
}
