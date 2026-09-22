package gotrail

type Option func(*tracerConfig)

func WithConnector(c Connector) Option {
	return func(t *tracerConfig) {
		t.connector = c
	}
}
