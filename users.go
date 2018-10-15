package fluidpay

import (
	"encoding/json"
)

// ChangePassword changes your password to the new one given in the request body.
func ChangePassword(fluidpay Fluidpay, reqBody ChangePasswordRequest) (GeneralResponse, error) {
	var response GeneralResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", "change-password"}
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

// CreateUser creates a new user.
func CreateUser(fluidpay Fluidpay, reqBody CreateUserRequest) (UserResponse, error) {
	var response UserResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user"}
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

// UpdateUser updates the user identified by the ID.
func UpdateUser(fluidpay Fluidpay, reqBody UpdateUserRequest, userID string) (UserResponse, error) {
	var response UserResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", userID}
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

// GetUser gets a specific user identified by the ID.
func GetUser(fluidpay Fluidpay, userID string) (UserResponse, error) {
	var response UserResponse
	param := []string{"user", userID}
	resBody, err := DoRequest(fluidpay, "GET", param, nil)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetCurrentUser gets current user.
func GetCurrentUser(fluidpay Fluidpay) (UserResponse, error) {
	var response UserResponse
	param := []string{"user"}
	resBody, err := DoRequest(fluidpay, "GET", param, nil)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetUsers gets all the users.
func GetUsers(fluidpay Fluidpay) (UsersResponse, error) {
	var response UsersResponse
	param := []string{"users"}
	resBody, err := DoRequest(fluidpay, "GET", param, nil)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// DeleteUser deletes the user identified by the id.
func DeleteUser(fluidpay Fluidpay, userID string) (GeneralResponse, error) {
	var response GeneralResponse
	param := []string{"user", userID}
	resBody, err := DoRequest(fluidpay, "DELETE", param, nil)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
