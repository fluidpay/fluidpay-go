package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var tokenUsrReq = CreateUserRequest{
	Username: "testtoken32",
	Name:     "test token merchant",
	Phone:    "6308885555",
	Email:    "info@website.com",
	Timezone: "CET",
	Password: "T@est12345678",
	Status:   "active",
	Role:     "admin",
}

var tokReq = JWTTokenRequest{
	Username: "testtoken32",
	Password: "T@est12345678",
}

var forUsrReq = ForgottenUsernameRequest{
	Email: "info@website.com",
}

var forPwReq = ForgottenPasswordRequest{
	Username: "testtoken32",
}

func TestAuth(t *testing.T) {
	auth, err := NewAuth("myAPIkey", 1)
	if err != nil {
		t.Error(err)
	}
	if auth.authorization != "myAPIkey" {
		t.Errorf("authorization should be 'myAPIkey', instead of '%v'", auth.authorization)
	}

	err = auth.SetAuth("myJWTtoken", 2)
	if err != nil {
		t.Error(err)
	}
	if auth.authorization != "Bearer myJWTtoken" {
		t.Errorf("authorization should be 'Bearer myJWTtoken', instead of '%v'", auth.authorization)
	}
}

func TestToken(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	tmpUsrRes, err := CreateUser(fluidpay, tokenUsrReq)
	if err != nil {
		t.Error(err)
	}

	tokRes, err := ObtainJWT(fluidpay, tokReq)
	if err != nil {
		t.Error(err)
	}
	if tokRes.Status != "success" {
		t.Errorf("Status should be 'success', instead of '%v'", tokRes.Status)
	}

	forUsrRes, err := ForgottenUsername(fluidpay, forUsrReq)
	if err != nil {
		t.Error(err)
	}
	if forUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", forUsrRes.Msg)
	}

	forPwRes, err := ForgottenPassword(fluidpay, forPwReq)
	if err != nil {
		t.Error(err)
	}
	if forPwRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", forPwRes.Msg)
	}

	/* err = TokenLogout(fluidpay)
	if err != nil {
		t.Error(err)
	} */

	DeleteUser(fluidpay, tmpUsrRes.Data.ID)
}
