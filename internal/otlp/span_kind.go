package otlp

type spanKind int

const (
	SpanKindUnspecified spanKind = iota
	SpanKindInternal
	SpanKindServer
	SpanKindClient
	SpanKindProducer
	SpanKindConsumer
)
