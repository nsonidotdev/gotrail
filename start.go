package tracer

import (
	"context"
	"time"

	"github.com/nsonidotdev/gotrail/internal/id"
)

func StartSpan(parent context.Context, name string, attrs map[string]any) context.Context {
	if !isInitialized.Load() || isTerminated.Load() {
		return parent
	}

	parentSpan, _ := getCtxSpan(parent)
	if parentSpan != nil && parentSpan.trace.isFinished.Load() {
		return parent
	}

	id, err := id.GenerateSpanID()
	if err != nil {
		return parent
	}

	newSpan := &Span{
		id:         id,
		name:       name,
		start:      time.Now(),
		attributes: attrs,
		status:     statusRunning,
		parent:     parentSpan,
	}

	if parentSpan == nil {
		newSpan.trace = newTrace(newSpan)
	} else {
		newSpan.trace = parentSpan.trace
	}

	tracer.tracker.recordStart(newSpan)

	if parentSpan != nil {
		appendSpanChild(parentSpan, newSpan)
	}

	ctx := setCtxSpan(parent, newSpan)
	stop := context.AfterFunc(ctx, func() {
		newSpan.finish(finishOptions{status: statusFail, reason: errSpanCtxCancelled.Error()})
	})
	newSpan.stopCancelListener = stop

	return ctx
}

func appendSpanChild(s *Span, child *Span) {
	if s == nil || child == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.children = append(s.children, child)
}
