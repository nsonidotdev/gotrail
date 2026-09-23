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
		if _, ok := span.Attributes["gotrail.skipped"]; !ok && span.status == statusSkip {
			span.Attributes["gotrail.skipped"] = true
		}
		attributes := otlp.SerializeAttributes(span.Attributes)
		status := toOTLPSpanStatus(span)
		otlpSpan := &otlp.Span{
			TraceID:           hex.EncodeToString(span.Trace.ID[:]),
			SpanID:            hex.EncodeToString(span.ID[:]),
			Name:              span.Name,
			StartTimeUnixNano: strconv.FormatInt(span.Start.UnixNano(), 10),
			EndTimeUnixNano:   strconv.FormatInt(span.Start.UnixNano()+int64(span.Duration), 10),
			Kind:              otlp.SpanKindServer,
			Status:            status,
			Attributes:        attributes,
		}
		if span.parent != nil {
			otlpSpan.ParentSpanID = hex.EncodeToString(span.parent.ID[:])
		}
		otlpSpans = append(otlpSpans, otlpSpan)
	}

	scopeAttrs := []*otlp.Attribute{}

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

func toOTLPSpanStatus(s *Span) *otlp.SpanStatus {
	switch s.status {
	case statusFail:
		return &otlp.SpanStatus{
			Code:    otlp.StatusError,
			Message: s.reason,
		}
	case statusSuccess:
		return &otlp.SpanStatus{
			Code: otlp.StatusOK,
		}

	case statusSkip:
		return &otlp.SpanStatus{
			Code:    otlp.StatusOK,
			Message: s.reason,
		}
	}

	return &otlp.SpanStatus{
		Code: otlp.StatusUnset,
	}
}
