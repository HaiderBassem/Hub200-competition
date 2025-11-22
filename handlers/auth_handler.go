package handlers

import (
	"googleforms/internal/dto"
	"googleforms/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login Login
// @Summary Login
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// check the data
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// check
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Email and password are required",
		})
	}

	// call service
	response, err := h.authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register user login
// @Summary user login
// @Description create new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "register hte data"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// check from the data
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Email, password, and full name are required",
		})
	}

	// call service
	user, err := h.authService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// GetCurrentUser get the data of the current user
// @Summary the data of the current user
// @Description get the data of the current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Data:    user,
	})
}
