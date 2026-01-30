package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func VerifyPayment(orderID, paymentID, signature string) bool {
	secret := os.Getenv("RAZORPAY_SECRET")

	data := orderID + "|" + paymentID

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generateSignature := hex.EncodeToString(h.Sum(nil))

	return generateSignature == signature
}
