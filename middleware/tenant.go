package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// TenantAccess middleware - access
func TenantAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userTenantID := c.Locals("tenantID")
		if userTenantID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "User not authenticated",
			})
		}

		// check the tenant
		tenantIDParam := c.Params("tenant_id")
		if tenantIDParam != "" {
			requestedTenantID, err := strconv.Atoi(tenantIDParam)
			if err != nil || uint(requestedTenantID) != userTenantID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error":   "Access denied to this tenant",
				})
			}
		}

		tenantIDQuery := c.Query("tenant_id")
		if tenantIDQuery != "" {
			requestedTenantID, err := strconv.Atoi(tenantIDQuery)
			if err != nil || uint(requestedTenantID) != userTenantID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error":   "Access denied to this tenant",
				})
			}
		}

		return c.Next()
	}
}

// check the context
func TenantRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenantID")
		if tenantID == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Tenant context missing",
			})
		}
		return c.Next()
	}
}
