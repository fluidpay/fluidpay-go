package fluidpay

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

// TestAPIKey is the api key used in the sandbox enviroment.
const TestAPIKey = "api_0wUsHIlrkK1I6ADno5MfT10UjhR"
const urlProduction = "https://app.fluidpay.com/api"
const urlSandbox = "https://sandbox.fluidpay.com/api"
const urlLocalDev = "http://localhost:8001/api"

// Fluidpay contains the api key, client and context.
type Fluidpay struct {
	APIKey   string
	Client   *http.Client
	Ctx      context.Context
	Sandbox  bool
	LocalDev bool
}

func urlBuilder(params []string) string {
	path := ""
	for _, param := range params {
		path += "/" + param
	}
	return path
}

// DoRequest returns a JSON from the API
// given the correct method, endpoint, JSON body (nil is accepted) and Fluidpay.
func DoRequest(fluidpay Fluidpay, method string, params []string, body []byte) ([]byte, error) {
	var res []byte
	var req *http.Request
	var err error
	var url string
	if fluidpay.LocalDev {
		url = urlLocalDev + urlBuilder(params)
	} else if fluidpay.Sandbox {
		url = urlSandbox + urlBuilder(params)
	} else {
		url = urlProduction + urlBuilder(params)
	}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return res, err
	}
	req.WithContext(fluidpay.Ctx)

	req.Header.Set("Authorization", fluidpay.APIKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := fluidpay.Client.Do(req)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	return res, err
}
