package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	ProductName string             `bson:"prod_name" json:"prod_name"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id" json:"product_id"`
	Product_Name *string            `json:"prod_name" bson:"prod_name"`
	Price        float64            `json:"price" bson:"price"`
	Rating       *uint              `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
	Quantity     *int               `json:"quantity" bson:"quantity"`
}
