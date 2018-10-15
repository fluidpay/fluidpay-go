# fluidpay-go
This is the official SDK for the Fluid Pay API written in the Go programming language.

## Getting started
After installing the SDK with ```go get github.org/fluidpay/fluidpay-go``` you first need to initiate the package.
```
import ( "fluidpay-go" )

fluidpay := Fluidpay{
		APIKey:   "yourApiKey",
		Client:   {http.Client},
		Ctx:      {context.Context},
	}
```

## General information
All communication with the API is struct based.
```
var creUsrReq = CreateUserRequest{
	Username: "yourUsername",
	Name:     "your name",
	Phone:    "6305555555",
	Email:    "your@email.address",
	Timezone: "CET",
	Password: "yourPassword",
	Status:   "active",
	Role:     "standard",
}
var creUsrRes UserResponse
var err error

creUsrRes, err = CreateUser(fluidpay, creUsrReq)
if err != nil {
	// handle error
}
```

## Documentation
Further information can be found on the [Fluid Pay website](https://sandbox.fluidpay.com/docs/index.html#introduction).