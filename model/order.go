package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID               primitive.ObjectID `json:"_id,omitempty"`
	Payment_Status   string             `json:"payment_status" bson:"payment_status"`
	Razorpay_OrderId string             `json:"razorpay_orderid" bson:"razorpay_orderid"`
	Total_price      float64            `json:"total_price" bson:"total_price"`
	Payment_Method   string             `json:"payment" bson:"payment"`
	Order_At         time.Time          `json:"time" bson:"time"`
	OrderCart        []ProductUser      `json:"ordercart" bson:"ordercart"`
	PaymentID        string             `json:"razorpay_paymentid" bson:"razorpay_paymentid"`
}
