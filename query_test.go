package fluidpay

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"
)

func TestQueryConstructors(t *testing.T) {
	equal(t, "Equal", *Equal("a"), StringQuery{Operator: "=", Value: "a"})
	equal(t, "NotEqual", *NotEqual("a"), StringQuery{Operator: "!=", Value: "a"})
	equal(t, "IntEqual", *IntEqual(1), IntQuery{Operator: "=", Value: 1})
	equal(t, "IntNotEqual", *IntNotEqual(1), IntQuery{Operator: "!=", Value: 1})
	equal(t, "LessThan", *LessThan(1), IntQuery{Operator: "<", Value: 1})
	equal(t, "GreaterThan", *GreaterThan(1), IntQuery{Operator: ">", Value: 1})

	chicago := time.FixedZone("CST", -6*3600)
	start := time.Date(2024, 3, 1, 18, 0, 0, 0, chicago)
	end := time.Date(2024, 3, 2, 18, 0, 0, 0, chicago)
	equal(t, "DateRange", *DateRange(start, end), DateRangeQuery{StartDate: "2024-03-02T00:00:00Z", EndDate: "2024-03-03T00:00:00Z"})
	equal(t, "Day", *Day(time.Date(2024, 3, 1, 23, 59, 0, 0, time.UTC)), DateRangeQuery{StartDate: "2024-03-01T00:00:00Z", EndDate: "2024-03-01T23:59:59Z"})
}

func TestQueryJSON(t *testing.T) {
	// The client's encoder must not HTML-escape the "<" operator.
	b, err := encodeJSON(TransactionSearchRequest{Amount: LessThan(500), CustomerID: Equal("c")})
	mustNoError(t, err)
	equal(t, "json", string(b), `{"amount":{"operator":"<","value":500},"customer_id":{"operator":"=","value":"c"}}`)

	b, err = json.Marshal(TransactionSearchRequest{})
	mustNoError(t, err)
	equal(t, "empty", string(b), `{}`)
}

func TestNewIdempotencyKey(t *testing.T) {
	uuidV4 := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		k := NewIdempotencyKey()
		if !uuidV4.MatchString(k) {
			t.Fatalf("%q is not a v4 UUID", k)
		}
		if seen[k] {
			t.Fatalf("duplicate key %q", k)
		}
		seen[k] = true
	}
}
