// Package otlp is a compatibility layer between the package
// internals and otlp `ExportTraceServiceRequest` expected payload
//
// You can see example trace payload here
// https://github.com/open-telemetry/opentelemetry-proto/blob/main/examples/trace.json
package otlp

import (
	"time"

	"github.com/nsonidotdev/gotrail/internal/id"
)

type ExportTraceServiceRequest struct {
	resourceSpans resourceSpan
}

type resourceSpan struct {
	resource   resource
	scopeSpans scopeSpan
}

type scopeSpan struct {
	scope
}

type resource struct {
	attributes []attr
}

type scope struct {
	name       string
	version    string
	attributes []attr
}

type span struct {
	traceID           id.TraceID
	spanID            id.SpanID
	parentSpanID      id.SpanID
	name              string
	startTimeUnixNano time.Duration
	endTimeUnixNano   time.Duration
	kind              spanKind
	attributes        []attr
}
