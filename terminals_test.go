package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

func TestTerminals(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	getTerRes, err := GetTerminals(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getTerRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", getTerRes.Status)
	}
	if getTerRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be '%v'", getTerRes.TotalCount)
	}
}
