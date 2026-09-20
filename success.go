package tracer

import "context"

func Success(ctx context.Context) {
	finish(ctx, finishOptions{
		status: statusSuccess,
	})
}

func (s *Span) Success() {
	s.finish(finishOptions{
		status: statusSuccess,
	})
}
