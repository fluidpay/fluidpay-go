package fluidpay

import (
	"encoding/json"
	"fmt"
)

// CreateCustomer creates a new customer token.
func CreateCustomer(fluidpay Fluidpay, reqBody CreateCustomerRequest) (CustomerResponse, error) {
	var response CustomerResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer"}
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

// GetCustomer returns a specific customer token identified by ID.
func GetCustomer(fluidpay Fluidpay, customerID string) (CustomerResponse, error) {
	var response CustomerResponse
	param := []string{"customer", customerID}
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

// UpdateCustomer updates a specific customer token identified by ID.
func UpdateCustomer(fluidpay Fluidpay, reqBody UpdateCustomerRequest, customerID string) (CustomerResponse, error) {
	var response CustomerResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer", customerID}
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

// DeleteCustomer deletes a specific customer token identified by ID.
func DeleteCustomer(fluidpay Fluidpay, customerID string) (CustomerResponse, error) {
	var response CustomerResponse
	param := []string{"customer", customerID}
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

// --------------------- Customer address section. --------------------------

// CreateCustomerAddress creates a customer address token.
func CreateCustomerAddress(fluidpay Fluidpay, reqBody CustomerAddressRequest, customerID string) (CustomerAddressResponse, error) {
	var response CustomerAddressResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer", customerID, "address"}
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

// GetCustomerAddresses gets all the address tokens of a customer identified by ID.
func GetCustomerAddresses(fluidpay Fluidpay, customerID string) (CustomerAddressesResponse, error) {
	var response CustomerAddressesResponse
	param := []string{"customer", customerID, "addresses"}
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

// GetCustomerAddress gets an address token of a customer both identified by ID.
func GetCustomerAddress(fluidpay Fluidpay, customerID, addressTokenID string) (CustomerAddressesResponse, error) {
	var response CustomerAddressesResponse
	param := []string{"customer", customerID, "address", addressTokenID}
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

// UpdateCustomerAddress updates an address token of a customer both identified by ID.
func UpdateCustomerAddress(fluidpay Fluidpay, reqBody CustomerAddressRequest, customerID, addressTokenID string) (CustomerAddressResponse, error) {
	var response CustomerAddressResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer", customerID, "address", addressTokenID}
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

// DeleteCustomerAddress deletes an address token of a customer both identified by ID.
func DeleteCustomerAddress(fluidpay Fluidpay, customerID, addressTokenID string) (CustomerAddressResponse, error) {
	var response CustomerAddressResponse
	param := []string{"customer", customerID, "address", addressTokenID}
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

// --------------------- Customer payment section. --------------------------

// CreateCustomerPayment creates a customer payment token.
func CreateCustomerPayment(fluidpay Fluidpay, reqBody CustomerPaymentRequest, customerID, paymentType string) (CustomerPaymentResponse, error) {
	var response CustomerPaymentResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer", customerID, "paymentmethod", paymentType}
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

// GetCustomerPayments gets all the payment tokens for a customer identified by ID.
func GetCustomerPayments(fluidpay Fluidpay, customerID, paymentType string) (CustomerPaymentsResponse, error) {
	var response CustomerPaymentsResponse
	param := []string{"customer", customerID, "paymentmethod", paymentType}
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

// GetCustomerPayment gets a payment token for a customer both identified by ID.
func GetCustomerPayment(fluidpay Fluidpay, customerID, paymentType, paymentTokenID string) (CustomerPaymentsResponse, error) {
	var response CustomerPaymentsResponse
	param := []string{"customer", customerID, "paymentmethod", paymentType, paymentTokenID}
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

// UpdateCustomerPayment updates a payment token of a customer both identified by ID.
func UpdateCustomerPayment(fluidpay Fluidpay, reqBody CustomerPaymentRequest, customerID, paymentType, paymentTokenID string) (CustomerPaymentResponse, error) {
	var response CustomerPaymentResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"customer", customerID, "paymentmethod", paymentType, paymentTokenID}
	resBody, err := DoRequest(fluidpay, "POST", param, body)
	fmt.Println(string(resBody))
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// DeleteCustomerPayment deletes a payment token of a customer both identified by ID.
func DeleteCustomerPayment(fluidpay Fluidpay, customerID, paymentType, paymentTokenID string) (CustomerPaymentResponse, error) {
	var response CustomerPaymentResponse
	param := []string{"customer", customerID, "paymentmethod", paymentType, paymentTokenID}
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
