package payment

import (
	"fmt"

	"github.com/razorpay/razorpay-go"
)

var razropay_client *razorpay.Client

func InitRazorpay(key, secret string) {
	razropay_client = razorpay.NewClient(key, secret)
}

func CreatePaymentOrder(amount int64) (string, error) {
	// 1. Check if client is initialized
	if razropay_client == nil {
		fmt.Println("❌ Razorpay Client not initialized!")
		return "", fmt.Errorf("razorpay client is nil")
	}

	data := map[string]interface{}{
		"amount":   amount, // Ensure this is (Actual Price * 100)
		"currency": "INR",
	}

	order, err := razropay_client.Order.Create(data, nil)
	if err != nil {
		// Isse terminal mein dekho, Razorpay kya error message de raha hai
		fmt.Printf("❌ Razorpay Error: %v\n", err)
		return "", err
	}

	orderID, ok := order["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get order_id from response")
	}

	return orderID, nil
}
