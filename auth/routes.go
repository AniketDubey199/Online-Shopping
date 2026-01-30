package auth

import (
	"context"
	"time"

	"github.com/AniketDubey199/online_shopping/database"
	"github.com/AniketDubey199/online_shopping/middleware"
	"github.com/AniketDubey199/online_shopping/model"
	"github.com/AniketDubey199/online_shopping/utils"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Usercollection *mongo.Collection = database.UserData(database.Client, "User")
var Prodcollection *mongo.Collection = database.ProdData(database.Client, "Product")

func Authentication(router fiber.Router) {
	router.Post("/register", func(c fiber.Ctx) error {
		user := &model.User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username or password cannot be empty",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		count, _ := Usercollection.CountDocuments(ctx, bson.M{"username": user.Username})

		if count > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Another username found",
			})
		}

		hashed, err := utils.Hashing(user.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password hashing cannot be done",
			})
		}

		user.Password = hashed

		_, err = Usercollection.InsertOne(ctx, user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "User registered Succesfully",
		})

	})
	router.Post("/login", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user := new(model.User)
		authUser := &model.User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}
		if authUser.Username == "" || authUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username or password cannot be empty",
			})
		}

		filter := bson.M{"username": authUser.Username}

		err := Usercollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "No Documents found",
				})
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}

		if err := utils.ComparePassword(user.Password, authUser.Password); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert",
			})
		}

		token, err := middleware.GenerateToken(user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannnot generate token",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"token": token,
		})
	})
}
