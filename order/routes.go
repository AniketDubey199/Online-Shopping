package order

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AniketDubey199/online_shopping/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProdToCart(ctx context.Context, useridcollection *mongo.Collection, productidcollection *mongo.Collection, prodid primitive.ObjectID, userid string, quantity int) error {
	var prod model.ProductUser
	err := productidcollection.FindOne(ctx, bson.M{"_id": prodid}).Decode(&prod)
	if err != nil {
		log.Printf("ERROR: Product DB mein nahi mila: %v", err)
		return err
	}

	if prod.Product_Name != nil {
		log.Printf("✅ SUCCESS: Product Picked! Name: %s, Price: %v", *prod.Product_Name, prod.Price)
	} else {
		log.Printf("⚠️ WARNING: Product picked but Name is NULL. Check your BSON tags!")
	}

	prod.Quantity = &quantity

	// 3. UserID ko ObjectID mein badlo
	uID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		log.Printf("ERROR: UserID conversion failed: %v", err)
		return err
	}

	// 4. Update Query (BSON tags check)
	filter := bson.M{"_id": uID}
	update := bson.M{
		"$push": bson.M{
			"usercart": prod,
		},
	}

	// 5. Execute Update
	_, err = useridcollection.UpdateOne(ctx, filter, update)
	if err != nil {
		// 3. AGAR ERROR AAYA: Iska matlab usercart 'null' hai
		log.Println("⚠️ Push fail hua (shyad cart null hai), initializing with $set...")

		// $set use karke null ko array se replace karo
		updateSet := bson.M{
			"$set": bson.M{
				"usercart": []model.ProductUser{prod},
			},
		}

		_, err = useridcollection.UpdateOne(ctx, filter, updateSet)
		if err != nil {
			log.Printf("❌ Set bhi fail ho gaya: %v", err)
			return fmt.Errorf("database error: cannot initialize cart")
		}

		log.Println("✅ Success: Cart initialized and product added!")
		return nil
	}

	log.Println("🔥 Success: Product pushed to existing cart!")
	return nil

}

func RemoveItemFromCart(ctx context.Context, useridcollection *mongo.Collection, productidcollection *mongo.Collection, prodid primitive.ObjectID, userid string) error {
	userID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userID}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": prodid}}}
	_, erro := useridcollection.UpdateOne(ctx, filter, update)
	if erro != nil {
		return erro
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, useridCollection *mongo.Collection, userID string) (primitive.ObjectID, int64, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}

	var user model.User
	err = useridCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}

	if len(user.UserCart) == 0 {
		log.Fatal("Cart is empty")
		return primitive.NilObjectID, 0, err
	}

	var total float64 = 0
	for _, item := range user.UserCart {
		total += item.Price * float64(*item.Quantity)
	}

	orderID := primitive.NewObjectID()

	orderdetails := model.Order{
		ID:             orderID,
		Order_At:       time.Now(),
		OrderCart:      user.UserCart,
		Total_price:    float64(total),
		Payment_Method: "Online",
		Payment_Status: "Pending",
	}

	update := bson.M{"$push": bson.M{"orders": orderdetails}}
	_, err = useridCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	return orderID, int64(total) * 100, nil
}

func CreateCODOrderFromCart(ctx context.Context, useridCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	var user model.User
	err = useridCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return err
	}

	if len(user.UserCart) == 0 {
		log.Fatal("Cart is empty")
		return err
	}

	var total float64 = 0
	for _, item := range user.UserCart {
		total += item.Price * float64(*item.Quantity)
	}

	orderdetails := model.Order{
		ID:             primitive.NewObjectID(),
		Order_At:       time.Now(),
		OrderCart:      user.UserCart,
		Total_price:    float64(total),
		Payment_Method: "COD",
		Payment_Status: "PAID",
	}

	update := bson.M{"$push": bson.M{"orders": orderdetails}}
	_, err = useridCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func InstantBuyer(ctx context.Context, useridcollection *mongo.Collection, productidcollection *mongo.Collection, prodid primitive.ObjectID, userid string) error {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}

	var product model.ProductUser
	ok := productidcollection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if ok != nil {
		return err
	}

	orderdetails := model.Order{
		ID:             primitive.NewObjectID(),
		Order_At:       time.Now(),
		OrderCart:      []model.ProductUser{product},
		Total_price:    float64(product.Price),
		Payment_Method: "Online",
		Payment_Status: "Pending",
	}

	update := bson.M{"$push": bson.M{"orders": orderdetails}}
	_, err = useridcollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}
