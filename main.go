package main

import (
	"log"
	"os"

	"github.com/AniketDubey199/online_shopping/auth"
	"github.com/AniketDubey199/online_shopping/cart"
	"github.com/AniketDubey199/online_shopping/database"
	"github.com/AniketDubey199/online_shopping/payment"
	"github.com/AniketDubey199/online_shopping/product"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	_ = godotenv.Load()

	app := fiber.New(fiber.Config{
		AppName: "E-Commerce",
	})

	app.Use(logger.New())

	db := database.InitializeDb()

	cartApp := cart.NewApplication(db.Collection("User"), db.Collection("Product"))

	auhtGroup := app.Group("/auth")

	auth.Authentication(auhtGroup, db)

	//public

	product.ProductGroup(app, db)

	//private

	cart.CartRoutes(app, cartApp)

	payment.InitRazorpay(
		os.Getenv("RAZORPAY_KEY"),
		os.Getenv("RAZORPAY_SECRET"),
	)

	log.Fatal(app.Listen("0.0.0.0:" + port))

}
