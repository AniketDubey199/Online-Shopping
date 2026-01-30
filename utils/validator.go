package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Authenticator() fiber.Handler {
	return func(c fiber.Ctx) error {
		var tokenstring string
		if cookiestring := c.Cookies("jwt"); cookiestring != "" {
			log.Printf("cookie token is being used")
			tokenstring = cookiestring
		} else {
			// retrieve it from the authenticator header
			autheader := c.Get("Authorization")
			if autheader == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Authorization is not in header",
				})
			}
			tokenparts := strings.Split(autheader, " ")
			if len(tokenparts) != 2 || strings.ToLower(tokenparts[0]) != "bearer" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Authrization is not valid",
				})
			}
			tokenstring = tokenparts[1]
		}
		if tokenstring == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No token provided",
			})
		}
		secret := []byte("super-secret-key")
		token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("Signing method not valid")
			}
			return secret, nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token", "details": err.Error()})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is not valid"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid claims has been requested",
			})
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No user id cannot be found",
			})
		}

		if _, err := primitive.ObjectIDFromHex(userID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid User Format",
			})
		}

		c.Locals("userID", userID)
		log.Printf("Middleware: UserID %s set in Locals", userID)

		return c.Next()
	}

}
