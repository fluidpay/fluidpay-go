package fluidpay

import (
	"encoding/json"
	"fmt"
)

// Auth is the struct to store your authentication.
type Auth struct {
	authType      int
	authorization string
}

// NewAuth initializes a new authentication.
// Set authType to 1 for API key and 2 for JWT token.
func NewAuth(auth string, authType int) (Auth, error) {
	switch authType {
	case 1:
		return Auth{
			authType:      authType,
			authorization: auth,
		}, nil
	case 2:
		return Auth{
			authType:      authType,
			authorization: "Bearer " + auth,
		}, nil
	default:
		return Auth{}, fmt.Errorf("Invalid authType %v please give 1 or 2 as a value", authType)
	}
}

// SetAuth updates the previous authentication.
// Set authType to 1 for API key and 2 for JWT token.
func (a *Auth) SetAuth(auth string, authType int) error {
	switch authType {
	case 1:
		a.authType = authType
		a.authorization = auth
		return nil
	case 2:
		a.authType = authType
		a.authorization = "Bearer " + auth
		return nil
	default:
		return fmt.Errorf("invalid authType %v please give 1 or 2 as a value", authType)
	}
}

// ObtainJWT obtains a new JWT token.
func ObtainJWT(fluidpay Fluidpay, reqBody JWTTokenRequest) (JWTTokenResponse, error) {
	var response JWTTokenResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"token-auth"}
	resBody, err := DoRequest(fluidpay, "POST", param, body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

//ForgottenUsername requests a reminder e-mail for the username.
func ForgottenUsername(fluidpay Fluidpay, reqBody ForgottenUsernameRequest) (GeneralResponse, error) {
	var response GeneralResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", "forgot-username"}
	resBody, err := DoRequest(fluidpay, "POST", param, body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

//ForgottenPassword requests a reminder e-mail containing a reset code for the password.
func ForgottenPassword(fluidpay Fluidpay, reqBody ForgottenPasswordRequest) (GeneralResponse, error) {
	var response GeneralResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", "forgot-password"}
	resBody, err := DoRequest(fluidpay, "POST", param, body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

//PasswordReset resets your password to the one given in the request body.
func PasswordReset(fluidpay Fluidpay, reqBody PasswordResetRequest) (GeneralResponse, error) {
	var response GeneralResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", "forgot-password", "reset"}
	resBody, err := DoRequest(fluidpay, "POST", param, body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

//TokenLogout terminates a valid authentication token.
func TokenLogout(fluidpay Fluidpay) error {
	param := []string{"logout"}
	_, err := DoRequest(fluidpay, "GET", param, nil)
	if err != nil {
		return err
	}
	return nil
}
