package fluidpay

import (
	"encoding/json"
)

// @TODO Sort out the TerminalResponse (need a valid terminal ID)

// DoTransaction handles one transaction process.
func DoTransaction(fluidpay Fluidpay, reqBody TransactionRequest) (TerminalResponse, error) {
	var response TerminalResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"transaction"}
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

// GetTransactionStatus gets the status of a transaction identified by ID.
func GetTransactionStatus(fluidpay Fluidpay, transactionID string) (TransactionSearchResponse, error) {
	var response TransactionSearchResponse
	param := []string{"transaction", transactionID}
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

// QueryTransactions gets transactions identified by the values in the request body.
func QueryTransactions(fluidpay Fluidpay, reqBody TransactionQueryRequest) (TransactionSearchResponse, error) {
	var response TransactionSearchResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"transaction", "search"}
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

// CaptureTransaction captures an already authorized transaction identified by ID.
func CaptureTransaction(fluidpay Fluidpay, reqBody TransactionCaptureRequest, transactionID string) (TransactionResponse, error) {
	var response TransactionResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"transaction", transactionID, "capture"}
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

// VoidTransaction voids a transaction with pending settlement identified by ID.
func VoidTransaction(fluidpay Fluidpay, transactionID string) (TransactionResponse, error) {
	var response TransactionResponse
	param := []string{"transaction", transactionID, "void"}
	resBody, err := DoRequest(fluidpay, "POST", param, nil)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// RefundTransaction refunds a previously settled amount.
func RefundTransaction(fluidpay Fluidpay, reqBody TransactionRefundRequest, transactionID string) (TransactionResponse, error) {
	var response TransactionResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"transaction", transactionID, "refund"}
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
