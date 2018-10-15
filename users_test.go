package fluidpay

import (
	"context"
	"net/http"
	"testing"
)

var creUsrReq = CreateUserRequest{
	Username: "testmerchantgo1",
	Name:     "test merchant user",
	Phone:    "6305555555",
	Email:    "info@website.com",
	Timezone: "CET",
	Password: "T@est12345678",
	Status:   "active",
	Role:     "admin",
}

var chPwReq = ChangePasswordRequest{
	Username:        "testmerchantgo1",
	CurrentPassword: "T@est12345678",
	NewPassword:     "T@est12345678",
}

var updUsrReq = UpdateUserRequest{
	Name:     "fresh test merchant user",
	Phone:    "6305555558",
	Email:    "info@site.com",
	Timezone: "CET",
	Status:   "active",
	Role:     "admin",
}

func TestUser(t *testing.T) {
	fluidpay := Fluidpay{
		APIKey:   TestAPIKey,
		Client:   new(http.Client),
		Ctx:      context.Background(),
		LocalDev: true,
	}

	currUsrRes, err := GetCurrentUser(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if currUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", currUsrRes.Msg)
	}

	getUsrsRes, err := GetUsers(fluidpay)
	if err != nil {
		t.Error(err)
	}
	if getUsrsRes.TotalCount == 0 {
		t.Errorf("TotalCount shouldn't be '%v'", getUsrsRes.TotalCount)
	}

	creUsrRes, err := CreateUser(fluidpay, creUsrReq)
	if err != nil {
		t.Error(err)
	}
	if creUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", creUsrRes.Msg)
	}

	id := creUsrRes.Data.ID

	getUsrRes, err := GetUser(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if getUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", getUsrRes.Msg)
	}

	/* chPwRes, err := ChangePassword(fluidpay, chPwReq)
	if err != nil {
		t.Error(err)
	}
	if chPwRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", chPwRes.Msg)
	} */

	updUsrRes, err := UpdateUser(fluidpay, updUsrReq, id)
	if err != nil {
		t.Error(err)
	}
	if updUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", updUsrRes.Msg)
	}

	delUsrRes, err := DeleteUser(fluidpay, id)
	if err != nil {
		t.Error(err)
	}
	if delUsrRes.Msg != "success" {
		t.Errorf("Msg should be 'success', instead of '%v'", delUsrRes.Msg)
	}
}
