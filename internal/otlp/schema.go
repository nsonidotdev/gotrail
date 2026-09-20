// Package otlp is a compatibility layer between the package
// internals and otlp `ExportTraceServiceRequest` expected payload
//
// You can see example trace payload here
// https://github.com/open-telemetry/opentelemetry-proto/blob/main/examples/trace.json
package otlp

import (
	"github.com/nsonidotdev/gotrail/internal/id"
)

type ExportTraceServiceRequest struct {
	ResourceSpans []ResourceSpan
}

type ResourceSpan struct {
	Resource   Resource
	ScopeSpans []ScopeSpan
}

type Resource struct {
	Attributes []Attribute
}

type ScopeSpan struct {
	Scope Scope
	Spans []Span
}

type Scope struct {
	Name       string
	Version    string
	Attributes []Attribute
}

type Span struct {
	TraceID           id.TraceID
	SpanID            id.SpanID
	ParentSpanID      id.SpanID
	Name              string
	StartTimeUnixNano uint64
	EndTimeUnixNano   uint64
	Kind              SpanKind
	Attributes        []Attribute
}
