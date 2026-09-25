package gotrail

import (
	"strconv"

	"github.com/nsonidotdev/gotrail/internal/otlp"
)

// Transforms spans to OTLP format
func SpansToOTLPJSON(spans []*Span) (*otlp.ExportTraceServiceRequest, error) {
	if len(spans) == 0 {
		return nil, errNoSpans
	}
	otlpSpans := make([]*otlp.Span, len(spans))
	for i, span := range spans {
		otlpSpans[i] = toOTLPSpan(span)
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

func toOTLPSpan(span *Span) *otlp.Span {
	span.mu.Lock()
	defer span.mu.Unlock()

	if _, ok := span.Attributes["skipped"]; !ok && span.status == statusSkip {
		span.Attributes["skipped"] = true
	}
	attributes := otlp.SerializeAttributes(span.Attributes)
	status := toOTLPSpanStatus(span)

	events := make([]*otlp.Event, len(span.Events))
	for i, event := range span.Events {
		events[i] = toOTLPEvent(event)
	}

	otlpSpan := &otlp.Span{
		TraceID:           span.Trace.ID.String(),
		SpanID:            span.ID.String(),
		Name:              span.Name,
		StartTimeUnixNano: strconv.FormatInt(span.Start.UnixNano(), 10),
		EndTimeUnixNano:   strconv.FormatInt(span.Start.UnixNano()+int64(span.Duration), 10),
		Kind:              otlp.SpanKindServer,
		Events:            events,
		Status:            status,
		Attributes:        attributes,
	}
	if span.parent != nil {
		otlpSpan.ParentSpanID = span.parent.ID.String()
	}

	return otlpSpan
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

func toOTLPEvent(e *Event) *otlp.Event {
	return &otlp.Event{
		Name:         e.Name,
		TimeUnixNano: strconv.FormatInt(e.StartTime.UnixNano(), 10),
		Attributes:   otlp.SerializeAttributes(e.Attributes),
	}
}
