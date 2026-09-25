package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/nsonidotdev/gotrail"
)

func main() {
	client := http.DefaultClient
	url := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	authorization := os.Getenv("OTEL_EXPORTER_OTLP_AUTHRIZATION_TOKEN")

	gotrail.Init(
		"server-stage",
		gotrail.WithConnector(
			gotrail.NewHTTPOTLPConnector(
				client,
				url,
				func(req *http.Request) error {
					req.Header.Set("Authorization", authorization)
					return nil
				},
			),
		),
	)

	rootCtx, root := gotrail.Start(context.Background(), "POST /api/checkout", map[string]any{
		"http.method": "POST",
		"http.route":  "/api/checkout",
		"user_id":     4821,
	})

	authCtx, auth := gotrail.Start(rootCtx, "authenticate", map[string]any{"strategy": "jwt"})
	time.Sleep(4 * time.Millisecond)
	_, verifyToken := gotrail.Start(authCtx, "verify-token", nil)
	time.Sleep(2 * time.Millisecond)
	verifyToken.AddEvent("token-verified", map[string]any{"exp_in_s": 3600})
	verifyToken.Success()
	auth.Success()

	cartCtx, cart := gotrail.Start(rootCtx, "load-cart", map[string]any{"cart_id": "c_88f2"})
	time.Sleep(6 * time.Millisecond)
	cart.AddEvent("cache-miss", map[string]any{"cache": "redis"})
	_, cartQuery := gotrail.Start(cartCtx, "db.query", map[string]any{
		"db.statement": "SELECT * FROM cart_items WHERE cart_id = $1",
	})
	time.Sleep(9 * time.Millisecond)
	cartQuery.Success()
	cart.Success()

	paymentCtx, payment := gotrail.Start(rootCtx, "charge-payment", map[string]any{
		"amount":   129.99,
		"currency": "USD",
	})
	time.Sleep(3 * time.Millisecond)
	_, gateway := gotrail.Start(paymentCtx, "stripe.charge", map[string]any{"provider": "stripe"})
	time.Sleep(20 * time.Millisecond)
	gateway.AddEvent("card-declined", map[string]any{"decline_code": "insufficient_funds"})
	gateway.Fail("card_declined")
	payment.Fail("payment gateway declined the charge")

	_, notify := gotrail.Start(rootCtx, "send-confirmation-email", map[string]any{"template": "order_confirmed"})
	time.Sleep(5 * time.Millisecond)
	notify.Skip("checkout failed, no confirmation to send")

	root.Fail("checkout aborted: payment declined")
}
