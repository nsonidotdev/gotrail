package gotrail

import (
	"encoding/json"
	"fmt"
)

// Entry point for handling completed spans
func handleTraceCompleted(t *Trace) {
	exportTraceRequestData, err := traceToOTLP(t)
	if err != nil {
		fmt.Println("error transforming trace to OTLP", err)
	}

	bodyBytes, err := json.Marshal(exportTraceRequestData)
	if err != nil {
		fmt.Println("marshallig error", err)
	}

	fmt.Println(string(bodyBytes))
}
