package fluidpay

import (
	"encoding/json"
)

//CreateKey creates a new API key for a user.
func CreateKey(fluidpay Fluidpay, reqBody KeyRequest) (KeyResponse, error) {
	var response KeyResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"user", "apikey"}
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

// GetKeys gets all the API keys to a user.
func GetKeys(fluidpay Fluidpay) (KeysResponse, error) {
	var response KeysResponse
	param := []string{"user", "apikeys"}
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

// DeleteKey deletes the API key identified by the id.
func DeleteKey(fluidpay Fluidpay, apiKeyID string) (KeyResponse, error) {
	var response KeyResponse
	param := []string{"user", "apikey", apiKeyID}
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
