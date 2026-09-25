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
		gotrail.WithConnector(
			gotrail.NewWriterConnector(file),
		),
	)

	checkoutApproved()
	getOrder()
	checkoutDeclined()
}

// checkoutApproved: full successful checkout.
func checkoutApproved() {
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
	gateway.AddEvent("charge-authorized", map[string]any{"auth_code": "a1b2c3"})
	gateway.Success()
	payment.Success()

	_, notify := gotrail.Start(rootCtx, "send-confirmation-email", map[string]any{"template": "order_confirmed"})
	time.Sleep(5 * time.Millisecond)
	notify.Success()

	root.Success()
}

// getOrder: simple read-path trace, no payment involved.
func getOrder() {
	rootCtx, root := gotrail.Start(context.Background(), "GET /api/orders/{id}", map[string]any{
		"http.method": "GET",
		"http.route":  "/api/orders/{id}",
		"order_id":    7734,
	})

	_, auth := gotrail.Start(rootCtx, "authenticate", map[string]any{"strategy": "jwt"})
	time.Sleep(3 * time.Millisecond)
	auth.Success()

	_, dbQuery := gotrail.Start(rootCtx, "db.query", map[string]any{
		"db.statement": "SELECT * FROM orders WHERE id = $1",
	})
	time.Sleep(8 * time.Millisecond)
	dbQuery.Success()

	root.Success()
}

// checkoutDeclined: payment gateway rejects the charge, so the checkout
// fails and the confirmation email step never runs.
func checkoutDeclined() {
	rootCtx, root := gotrail.Start(context.Background(), "POST /api/checkout", map[string]any{
		"http.method": "POST",
		"http.route":  "/api/checkout",
		"user_id":     5310,
	})

	authCtx, auth := gotrail.Start(rootCtx, "authenticate", map[string]any{"strategy": "jwt"})
	time.Sleep(4 * time.Millisecond)
	_, verifyToken := gotrail.Start(authCtx, "verify-token", nil)
	time.Sleep(2 * time.Millisecond)
	verifyToken.Success()
	auth.Success()

	cartCtx, cart := gotrail.Start(rootCtx, "load-cart", map[string]any{"cart_id": "c_91a7"})
	time.Sleep(6 * time.Millisecond)
	_, cartQuery := gotrail.Start(cartCtx, "db.query", map[string]any{
		"db.statement": "SELECT * FROM cart_items WHERE cart_id = $1",
	})
	time.Sleep(9 * time.Millisecond)
	cartQuery.Success()
	cart.Success()

	paymentCtx, payment := gotrail.Start(rootCtx, "charge-payment", map[string]any{
		"amount":   249.50,
		"currency": "USD",
	})
	time.Sleep(3 * time.Millisecond)
	_, gateway := gotrail.Start(paymentCtx, "stripe.charge", map[string]any{"provider": "stripe"})
	time.Sleep(15 * time.Millisecond)
	gateway.AddEvent("card-declined", map[string]any{"decline_code": "insufficient_funds"})
	gateway.Fail("card_declined")
	payment.Fail("payment gateway declined the charge")

	root.Fail("checkout aborted: payment declined")
}
