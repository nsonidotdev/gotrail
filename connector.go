package gotrail

type Connector interface {
	Send([]*Span) error
}
