# gotrail

A small Go library for building trace trees (spans with parent/child
relationships) in-process and shipping them out as OTLP.

This is a **learning project** — a way to understand how tracing / OTLP
works by building it from scratch, not a hardened observability SDK. It has
no performance guarantees and isn't optimized for production workloads. That
said, it's plain Go with no required third-party dependencies, so feel free
to use it if it fits your needs.

## Features

- Context-based spans: `Start` returns a child `context.Context` carrying the
  new span, so nested calls just pass `ctx` down.
- Explicit outcomes: every span ends with `Success()`, `Fail(reason)`, or
  `Skip(reason)`.
- Automatic parent/child nesting and per-trace completion tracking — a trace
  is flushed once its root and every descendant span has finished.
- Span cancellation: if the span's `context.Context` is canceled, the span is
  automatically failed.
- Pluggable output via the `Connector` interface — ship spans wherever you
  want.
- Two connectors included:
  - `WriterConnector` — writes OTLP JSON to any `io.Writer` (a file, stdout, etc.).
  - `HTTPOTLPConnector` — POSTs OTLP JSON to an HTTP OTLP endpoint (e.g. Grafana Tempo/Cloud, Honeycomb, any OTLP/HTTP collector).
- `PrintTrace` / `PrintSpanTree` for a quick human-readable tree of a trace
  in your terminal.

## Installation

```bash
go get github.com/nsonidotdev/gotrail
```

## Usage

Call `gotrail.Init` once at startup with a service name and a connector.
Then use `gotrail.Start` to open spans, threading the returned context
through your call chain, and end each span with `Success`, `Fail`, or `Skip`.

### Writer connector

Writes OTLP JSON to any `io.Writer` — useful for local development or
writing traces to a file.

```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/nsonidotdev/gotrail"
)

func main() {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gotrail.Init(
		"server-stage",
		gotrail.WithConnector(gotrail.NewWriterConnector(file)),
	)

	rootCtx, root := gotrail.Start(context.Background(), "POST /api/checkout", map[string]any{
		"http.method": "POST",
		"http.route":  "/api/checkout",
		"user_id":     4821,
	})

	_, auth := gotrail.Start(rootCtx, "authenticate", map[string]any{"strategy": "jwt"})
	time.Sleep(4 * time.Millisecond)
	auth.Success()

	root.Success()
}
```

### HTTP OTLP connector

POSTs OTLP JSON to any OTLP/HTTP endpoint (Grafana Cloud, Tempo, Honeycomb,
your own collector, ...). `PrepareRequest` lets you set auth headers before
each send.

```go
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/nsonidotdev/gotrail"
)

func main() {
	gotrail.Init(
		"server-stage",
		gotrail.WithConnector(
			gotrail.NewHTTPOTLPConnector(
				http.DefaultClient,
				os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
				func(req *http.Request) error {
					req.Header.Set("Authorization", os.Getenv("OTEL_EXPORTER_OTLP_AUTHRIZATION_TOKEN"))
					return nil
				},
			),
		),
	)

	rootCtx, root := gotrail.Start(context.Background(), "POST /api/checkout", map[string]any{
		"http.method": "POST",
		"http.route":  "/api/checkout",
	})

	_, payment := gotrail.Start(rootCtx, "charge-payment", map[string]any{"amount": 129.99})
	time.Sleep(20 * time.Millisecond)
	payment.Fail("card_declined")

	root.Fail("checkout aborted: payment declined")
}
```

More complete, runnable examples live in [`examples/writer`](examples/writer)
and [`examples/http_otlp`](examples/http_otlp).

### Writing your own connector

`Connector` is a single-method interface, so shipping spans anywhere else
(stdout, a queue, a different wire format) just means implementing it:

```go
type Connector interface {
	Send([]*Span) error
}
```

## Status

Learning project, API may change. No performance or stability guarantees —
use it, fork it, or read it to understand how OTLP tracing works under the
hood.
