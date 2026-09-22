package gotrail

import "fmt"

// Entry point for handling completed traces
func handleTraceCompleted(t *Trace) {
	PrintTrace(t)

	if tracer.connector == nil {
		return
	}

	flatSpans := flattenSpans(t)
	err := tracer.connector.Send(flatSpans)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
