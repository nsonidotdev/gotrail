package gotrail

import (
	"encoding/hex"
	"strconv"

	"github.com/nsonidotdev/gotrail/internal/otlp"
)

// Transforms spans to OTLP format
func SpansToOTLPJSON(spans []*Span) (*otlp.ExportTraceServiceRequest, error) {
	if len(spans) == 0 {
		return nil, errNoSpans
	}
	otlpSpans := make([]*otlp.Span, 0, len(spans))
	for _, span := range spans {
		attributes := otlp.SerializeAttributes(span.Attributes)
		otlpSpan := &otlp.Span{
			TraceID:           hex.EncodeToString(span.Trace.ID[:]),
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
		otlpSpans = append(otlpSpans, otlpSpan)
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
		Spans: otlpSpans,
	}
	scopeSpans := make([]*otlp.ScopeSpan, 0, 1)
	scopeSpans = append(scopeSpans, scopeSpan)

	resource := &otlp.Resource{
		Attributes: otlp.SerializeAttributes(map[string]any{
			"service.name": tracer.service,
		}),
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
