package id

import (
	"crypto/rand"
	"encoding/hex"
)

// OTel requirements for ids
type (
	TraceID [16]byte
	SpanID  [8]byte
)

func GenerateTraceID() (TraceID, error) {
	var id TraceID
	if _, err := rand.Read(id[:]); err != nil {
		return TraceID{}, err
	}

	return id, nil
}

func GenerateSpanID() (SpanID, error) {
	var id SpanID
	if _, err := rand.Read(id[:]); err != nil {
		return SpanID{}, err
	}

	return id, nil
}

func (id TraceID) String() string {
	return hex.EncodeToString(id[:])
}

func (id SpanID) String() string {
	return hex.EncodeToString(id[:])
}
