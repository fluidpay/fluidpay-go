package fluidpay

import (
	"encoding/json"
)

// --------------------- Add on section. --------------------------

// CreateAddOn creates a new add on.
func CreateAddOn(fluidpay Fluidpay, reqBody RecurrenceRequest) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "addon"}
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

// GetAddOn gets an add on identified by ID.
func GetAddOn(fluidpay Fluidpay, addOnID string) (RecurrencesResponse, error) {
	var response RecurrencesResponse
	param := []string{"recurring", "addon", addOnID}
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

// GetAddOns gets all add ons.
func GetAddOns(fluidpay Fluidpay) (RecurrencesResponse, error) {
	var response RecurrencesResponse
	param := []string{"recurring", "addons"}
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

// UpdateAddOn updates an add on identified by ID.
func UpdateAddOn(fluidpay Fluidpay, reqBody RecurrenceRequest, addOnID string) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "addon", addOnID}
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

// DeleteAddOn deletes an add on identified by ID.
func DeleteAddOn(fluidpay Fluidpay, addOnID string) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	param := []string{"recurring", "addon", addOnID}
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

// --------------------- Discount section. --------------------------

// CreateDiscount creates a new discount.
func CreateDiscount(fluidpay Fluidpay, reqBody RecurrenceRequest) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "discount"}
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

// GetDiscount gets a discount identified by ID.
func GetDiscount(fluidpay Fluidpay, discountID string) (RecurrencesResponse, error) {
	var response RecurrencesResponse
	param := []string{"recurring", "discount", discountID}
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

// GetDiscounts gets all the discounts.
func GetDiscounts(fluidpay Fluidpay) (RecurrencesResponse, error) {
	var response RecurrencesResponse
	param := []string{"recurring", "discounts"}
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

// UpdateDiscount updates a discount identified by ID.
func UpdateDiscount(fluidpay Fluidpay, reqBody RecurrenceRequest, discountID string) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "discount", discountID}
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

// DeleteDiscount deletes a discount identified by ID.
func DeleteDiscount(fluidpay Fluidpay, discountID string) (RecurrenceResponse, error) {
	var response RecurrenceResponse
	param := []string{"recurring", "discount", discountID}
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

// --------------------- Plan section. --------------------------

// CreatePlan creates a new plan.
func CreatePlan(fluidpay Fluidpay, reqBody PlanRequest) (PlanResponse, error) {
	var response PlanResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "plan"}
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

// UpdatePlan updates a plan identified by ID.
func UpdatePlan(fluidpay Fluidpay, reqBody PlanRequest, planID string) (PlanResponse, error) {
	var response PlanResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "plan", planID}
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

// GetPlan gets a plan identified by ID.
func GetPlan(fluidpay Fluidpay, planID string) (PlansResponse, error) {
	var response PlansResponse
	param := []string{"recurring", "plan", planID}
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

// GetPlans gets all the plans.
func GetPlans(fluidpay Fluidpay) (PlansResponse, error) {
	var response PlansResponse
	param := []string{"recurring", "plans"}
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

// DeletePlan deletes a plan identified by ID.
func DeletePlan(fluidpay Fluidpay, planID string) (PlanResponse, error) {
	var response PlanResponse
	param := []string{"recurring", "plan", planID}
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

// --------------------- Subscription section. --------------------------

// CreateSubscription creates a subscription.
func CreateSubscription(fluidpay Fluidpay, reqBody SubscriptionRequest) (SubscriptionResponse, error) {
	var response SubscriptionResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "subscription"}
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

// UpdateSubscription updates a subscription identified by ID.
func UpdateSubscription(fluidpay Fluidpay, reqBody SubscriptionRequest, subID string) (SubscriptionResponse, error) {
	var response SubscriptionResponse
	body, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}
	param := []string{"recurring", "subscription", subID}
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

// GetSubscription gets a subscription identified by ID.
func GetSubscription(fluidpay Fluidpay, subID string) (SubscriptionsResponse, error) {
	var response SubscriptionsResponse
	param := []string{"recurring", "subscription", subID}
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

// DeleteSubscription deletes a subscription identified by ID.
func DeleteSubscription(fluidpay Fluidpay, subID string) (SubscriptionResponse, error) {
	var response SubscriptionResponse
	param := []string{"recurring", "subscription", subID}
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
