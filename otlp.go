package gotrail

import (
	"encoding/hex"
	"strconv"

	"github.com/nsonidotdev/gotrail/internal/otlp"
)

// Transforms trace and all spans to OTLP format
func traceToOTLP(t *Trace) (*otlp.ExportTraceServiceRequest, error) {
	if t.root == nil {
		return nil, errNoSpans
	}
	flatSpans := flattenSpans(t)
	spans := make([]*otlp.Span, 0, len(flatSpans))
	for _, span := range flatSpans {
		attributes := otlp.SerializeAttributes(span.Attributes)
		otlpSpan := &otlp.Span{
			TraceID:           hex.EncodeToString(t.ID[:]),
			SpanID:            hex.EncodeToString(span.ID[:]),
			Name:              span.Name,
			StartTimeUnixNano: strconv.FormatInt(span.Start.UnixNano(), 10),
			EndTimeUnixNano:   strconv.FormatInt(span.Start.UnixNano()+int64(span.Duration), 10),
			Attributes:        attributes,
			Kind:              otlp.SpanKindServer,
		}
		if span.parent != nil {
			otlpSpan.ParentSpanID = hex.EncodeToString(span.parent.ID[:])
		}
		spans = append(spans, otlpSpan)
	}

	rawAttrs := map[string]any{
		"repo":       "https://github.com/nsonidotdev/gotrail",
		"suggestion": "star this repo please",
	}
	scopeAttrs := otlp.SerializeAttributes(rawAttrs)

	scope := &otlp.Scope{
		Name:       "gotrail",
		Version:    "v0.0.1",
		Attributes: scopeAttrs,
	}

	scopeSpan := &otlp.ScopeSpan{
		Scope: scope,
		Spans: spans,
	}
	scopeSpans := make([]*otlp.ScopeSpan, 0, 1)
	scopeSpans = append(scopeSpans, scopeSpan)

	resource := &otlp.Resource{
		// Attributes
	}

	resourceSpan := &otlp.ResourceSpan{
		Resource:   resource,
		ScopeSpans: scopeSpans,
	}
	resourceSpans := make([]*otlp.ResourceSpan, 0, 1)
	resourceSpans = append(resourceSpans, resourceSpan)

	return &otlp.ExportTraceServiceRequest{
		ResourceSpans: resourceSpans,
	}, nil
}
