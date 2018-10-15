package fluidpay

import (
	"encoding/json"
)

// GetTerminals gets all the terminals.
func GetTerminals(fluidpay Fluidpay) (TerminalsResponse, error) {
	var response TerminalsResponse
	param := []string{"terminals"}
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
