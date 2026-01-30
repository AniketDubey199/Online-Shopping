package product

import (
	"context"
	"time"

	"github.com/AniketDubey199/online_shopping/database"
	"github.com/AniketDubey199/online_shopping/model"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProductCollection *mongo.Collection = database.ProdData(database.Client, "Product")
var UserCollection *mongo.Collection = database.UserData(database.Client, "User")

func ProductGroup(app *fiber.App) {

	productGroup := app.Group("/products")

	productGroup.Get("/", SearchProduct)
	productGroup.Get("/search", SearchQueryProduct)
	productGroup.Post("/add", AddProduct)
}

func SearchProduct(c fiber.Ctx) error {
	var productlist []model.Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := ProductCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot retain productlist",
		})
	}

	err = cursor.All(ctx, &productlist)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot convert",
		})
	}

	defer cursor.Close(ctx)

	defer cancel()

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"product": productlist,
	})
}

func SearchQueryProduct(c fiber.Ctx) error {
	var searchproduct []model.Product
	queryParam := c.Query("name")

	if queryParam == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "No product is there",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	querySearch, err := ProductCollection.Find(ctx, bson.M{"prod_name": bson.M{"$regex": queryParam, "$options": "i"}})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot find it in the query",
		})
	}
	err = querySearch.All(ctx, &searchproduct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid",
		})
	}

	if len(searchproduct) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no product found",
		})
	}

	defer querySearch.Close(ctx)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"count":   len(searchproduct),
		"product": searchproduct,
	})

}

func AddProduct(c fiber.Ctx) error {

	var product model.Product

	if err := c.Bind().Body(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if product.ProductName == "" || product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product data",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := ProductCollection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot add product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "product added successfully",
		"productId": result.InsertedID,
	})
}
