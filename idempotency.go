package fluidpay

import (
	"crypto/rand"
	"fmt"
)

// NewIdempotencyKey returns a random UUID (version 4) suitable for
// TransactionRequest.IdempotencyKey. Reusing the same key within the
// idempotency window makes the gateway return the original transaction
// instead of charging the customer twice, which makes retries safe.
func NewIdempotencyKey() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failing is unrecoverable; surface it loudly rather
		// than silently returning a weak key.
		panic("fluidpay: crypto/rand unavailable: " + err.Error())
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
