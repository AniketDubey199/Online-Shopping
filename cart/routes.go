package cart

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/AniketDubey199/online_shopping/model"
	"github.com/AniketDubey199/online_shopping/order"
	"github.com/AniketDubey199/online_shopping/payment"
	"github.com/AniketDubey199/online_shopping/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CartRoutes(app *fiber.App, CartApp *Application) {

	protected := app.Group("/cart", utils.Authenticator())

	protected.Post("/add", CartApp.AddToCart)
	protected.Delete("/remove", CartApp.RemoveItem)
	protected.Get("/getitem", CartApp.GetItem)
	protected.Post("/buy", CartApp.BuyItemCOD)
	protected.Post("/instant-buy", CartApp.InstantBuy)
	protected.Post("/checkitem", CartApp.CheckoutCart)
	protected.Post("/verify", CartApp.VerifyPayment)

}

type Application struct {
	UserIDCollection    *mongo.Collection
	ProductIDCollection *mongo.Collection
}

func NewApplication(UserIDCollection, ProductIDCollection *mongo.Collection) *Application {
	return &Application{
		UserIDCollection:    UserIDCollection,
		ProductIDCollection: ProductIDCollection,
	}
}

func (app *Application) AddToCart(c fiber.Ctx) error {
	productQueryID := c.Query("id")
	if productQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is empty",
		})
	}
	quantityitem := c.Query("quantity", "1")

	quantity, _ := strconv.Atoi(quantityitem)
	if quantity <= 0 {
		quantity = 1
	}

	idRaw := c.Locals("userID")
	if idRaw == nil {
		log.Println("ERROR: Locals mein userID nahi mili!")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userQueryID := c.Locals("userID").(string)
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is empty",
		})
	}

	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product id",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = order.AddProdToCart(ctx, app.UserIDCollection, app.ProductIDCollection, productID, userQueryID, quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order cannot be added",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Item added successfully",
	})

}

func (app *Application) RemoveItem(c fiber.Ctx) error {
	productQueryID := c.Query("id")
	if productQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is empty",
		})
	}

	userQueryID := c.Locals("userID").(string) // ye karna zaroori hai kyunki ham jwt se userid lenge na ki database se
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is empty",
		})
	}

	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product id",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = order.RemoveItemFromCart(ctx, app.UserIDCollection, app.ProductIDCollection, productID, userQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot remove from cart",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Item removed successfully",
	})
}

func (app *Application) GetItem(c fiber.Ctx) error {
	user_id := c.Locals("userID").(string)
	if user_id == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "User Not found",
		})
	}

	usert_id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Cannot convert",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filledCart model.User
	err = app.UserIDCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Cannot Extract",
		})
	}

	filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}}}}}
	pointcursor, err := app.UserIDCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Some in pointer cursor",
		})
	}

	var listing []bson.M
	if err = pointcursor.All(ctx, &listing); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Cannot use it in listing",
		})
	}

	if len(listing) == 0 {
		return c.Status(200).JSON(fiber.Map{"total": 0, "message": "Cart is empty"})
	}

	result := listing[0]
	totalBill := result["total"]

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"total": totalBill,
	})
}

func (app *Application) BuyItemCOD(c fiber.Ctx) error {
	userQueryID := c.Locals("userID").(string)
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is empty",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := order.CreateCODOrderFromCart(ctx, app.UserIDCollection, userQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot Buy from cart",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Item Bought successfully",
	})

}

func (app *Application) CheckoutCart(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is empty",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderID, amount, err := order.BuyItemFromCart(ctx, app.UserIDCollection, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "There is some issue in fetching the cart payment",
		})
	}

	razorpaymentID, err := payment.CreatePaymentOrder(amount)
	if err != nil {
		return c.SendStatus(500)
	}

	_, err = app.UserIDCollection.UpdateOne(ctx, bson.M{"orders.id": orderID}, bson.M{"$set": bson.M{"orders.$.razorpay_orderid": razorpaymentID}})
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"order_id":         orderID,
		"razorpay_orderid": razorpaymentID,
		"amount":           amount,
	})
}

func (app *Application) VerifyPayment(c fiber.Ctx) error {
	var body struct {
		OrderID   string `json:"razorpay_orderid"`
		PaymentID string `json:"razorpay_paymentid"`
		Signature string `json:"razorpay_signature"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	isValid := payment.VerifyPayment(
		body.OrderID,
		body.PaymentID,
		body.Signature,
	)

	if !isValid {
		_, _ = app.UserIDCollection.UpdateOne(context.Background(),
			bson.M{"orders.razorpay_orderid": body.OrderID},
			bson.M{"$set": bson.M{"orders.$.payment_status": "FAILED"}},
		)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment verificatioin failed",
		})
	}

	_, err := app.UserIDCollection.UpdateOne(
		context.TODO(),
		bson.M{"orders.razorpay_orderid": body.OrderID},
		bson.M{"$set": bson.M{
			"orders.$.payment_status":     "PAID",
			"orders.$.razorpay_paymentid": body.PaymentID,
		}},
	)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "payment verified succefully",
	})
}

func (app *Application) InstantBuy(c fiber.Ctx) error {
	productQueryID := c.Query("id")
	if productQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is empty",
		})
	}

	userQueryID := c.Locals("userID").(string)
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is empty",
		})
	}

	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product id",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = order.InstantBuyer(ctx, app.UserIDCollection, app.ProductIDCollection, productID, userQueryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot Buy from cart",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Item Bought successfully",
	})

}
