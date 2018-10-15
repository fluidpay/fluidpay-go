package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var keyReq = KeyRequest{
	Type: "api",
	Name: "testapikey",
}

var keyUsrReq = CreateUserRequest{
	Username: "testkeygo001",
	Name:     "test key merchant",
	Phone:    "6308886655",
	Email:    "info@website.com",
	Timezone: "CET",
	Password: "T@est12345678",
	Status:   "active",
	Role:     "admin",
}

func TestKey(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	/* tmpUsrRes, err := CreateUser(fluidpay, keyUsrReq)
	if err != nil {
		t.Error(err)
	} */

	creKeyRes, err := CreateKey(fluidpay, keyReq)
	if err != nil {
		t.Error(err)
	}
	if creKeyRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creKeyRes.Msg)
	}

	id := creKeyRes.Data.APIKey

	getKeysRes, err := GetKeys(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getKeysRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be %v", getKeysRes.TotalCount)
	}

	delKeyRes, err := DeleteKey(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if delKeyRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delKeyRes.Msg)
	}

	//DeleteUser(fluidpay, id)
}
