// Command webhook runs an HTTP server that verifies and decodes FluidPay
// webhook deliveries.
//
//	FLUIDPAY_WEBHOOK_SECRET=<signing secret> go run ./examples/webhook
//
// Point a sandbox webhook at http://<host>:8080/webhooks/fluidpay and use the
// control panel's test button: the test delivery is verified with the shared
// test key automatically.
package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/fluidpay/fluidpay-go"
)

func main() {
	secret := os.Getenv("FLUIDPAY_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("FLUIDPAY_WEBHOOK_SECRET is required")
	}
	// The gateway may deliver the same event more than once; remember what
	// has been processed. Use a durable store in production.
	var (
		mu   sync.Mutex
		seen = map[string]bool{}
	)

	http.HandleFunc("/webhooks/fluidpay", func(w http.ResponseWriter, r *http.Request) {
		event, err := fluidpay.ParseWebhookRequest(r, secret)
		if err != nil {
			log.Printf("rejected webhook: %v", err)
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
		if event.IsTest() {
			log.Printf("test delivery from %s verified", event.AccountTypeID)
			w.WriteHeader(http.StatusOK)
			return
		}

		mu.Lock()
		duplicate := seen[event.EventID]
		seen[event.EventID] = true
		mu.Unlock()
		if duplicate {
			w.WriteHeader(http.StatusOK) // already handled; acknowledge again
			return
		}

		switch event.Type {
		case fluidpay.WebhookTypeTransactionCreate:
			tx, err := event.Transaction()
			if err != nil {
				log.Printf("event %s: %v", event.EventID, err)
				break
			}
			if tx.TransactionSource == fluidpay.SourceRecurring {
				log.Printf("subscription %s renewal %s: %s", tx.SubscriptionID, tx.ID, tx.Status)
			} else {
				log.Printf("transaction %s: %s (%d cents)", tx.ID, tx.Status, tx.Amount)
			}
		default:
			log.Printf("event %s of type %s", event.EventID, event.Type)
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
