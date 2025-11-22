package middleware

import (
	"googleforms/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Authorization token required",
			})
		}

		tokenString := strings.TrimSpace(authHeader)
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		user, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		c.Locals("user", user)
		c.Locals("userID", user.ID)
		c.Locals("tenantID", user.TenantID)
		c.Locals("userRole", user.Role)

		return c.Next()
	}
}

// AdminRequired middleware
func AdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "User not authenticated",
			})
		}

		role := userRole.(string)
		if role != "owner" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Admin access required",
			})
		}

		return c.Next()
	}
}

// OwnerRequired middleware
func OwnerRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "User not authenticated",
			})
		}

		role := userRole.(string)
		if role != "owner" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Owner access required",
			})
		}

		return c.Next()
	}
}
