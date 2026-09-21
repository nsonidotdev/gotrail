// Package otlp is a compatibility layer between the package
// internals and otlp `ExportTraceServiceRequest` expected payload
//
// You can see example trace payload here
// https://github.com/open-telemetry/opentelemetry-proto/blob/main/examples/trace.json
package otlp

type ExportTraceServiceRequest struct {
	ResourceSpans []*ResourceSpan `json:"resourceSpans"`
}

type ResourceSpan struct {
	Resource   *Resource    `json:"resource"`
	ScopeSpans []*ScopeSpan `json:"scopeSpans"`
}

type Resource struct {
	Attributes []*Attribute `json:"attributes,omitempty"`
}

type ScopeSpan struct {
	Scope *Scope  `json:"scope"`
	Spans []*Span `json:"spans"`
}

type Scope struct {
	Name       string       `json:"name"`
	Version    string       `json:"version"`
	Attributes []*Attribute `json:"attributes,omitempty"`
}

type Span struct {
	TraceID           string       `json:"traceId"`                // id.TraceID
	SpanID            string       `json:"spanId"`                 // id.SpanID
	ParentSpanID      string       `json:"parentSpanId,omitempty"` // id.SpanID
	Name              string       `json:"name"`
	StartTimeUnixNano string       `json:"startTimeUnixNano"`
	EndTimeUnixNano   string       `json:"endTimeUnixNano"`
	Kind              SpanKind     `json:"kind"`
	Attributes        []*Attribute `json:"attributes,omitempty"`
}
