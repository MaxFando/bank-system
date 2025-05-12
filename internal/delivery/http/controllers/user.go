package controllers

import (
	"github.com/MaxFando/bank-system/internal/core/bank/service/user"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *user.Service
	authService *user.AuthService
}

func NewUserController(userService *user.Service, authService *user.AuthService) *UserController {
	return &UserController{
		userService: userService,
		authService: authService,
	}
}

func (ctrl *UserController) Register(c echo.Context) error {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	id, err := ctrl.authService.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Registration failed"})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Registration successful",
		"user": map[string]interface{}{
			"id":    id,
			"email": req.Email,
		},
	})
}

func (ctrl *UserController) Login(c echo.Context) error {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	token, err := ctrl.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	// Example response
	return c.JSON(200, map[string]string{"message": "Login successful", "token": token})
}
