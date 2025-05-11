package http

import (
	"context"
	"github.com/MaxFando/bank-system/internal/providers"
	"github.com/MaxFando/bank-system/pkg/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewHandler(ctx context.Context, provider *providers.ServiceProvider) *echo.Echo {
	echoMainServer := echo.New()
	echoMainServer.Validator = &customValidator{validator: validator.New()}

	ctrl := NewController(provider)

	echoMainServer.POST("/register", ctrl.Register)
	echoMainServer.POST("/login", ctrl.Login)
	echoMainServer.POST("/accounts", ctrl.CreateAccount, echo.WrapMiddleware(auth.AuthMiddleware))
	echoMainServer.POST("/cards", ctrl.CreateCard, echo.WrapMiddleware(auth.AuthMiddleware))
	echoMainServer.POST("/transfer", ctrl.Transfer, echo.WrapMiddleware(auth.AuthMiddleware))
	echoMainServer.GET("/credits/{creditId}/schedule", ctrl.GetCreditSchedule, echo.WrapMiddleware(auth.AuthMiddleware))

	return echoMainServer
}

type Controller struct {
	serviceProvider *providers.ServiceProvider
}

func NewController(serviceProvider *providers.ServiceProvider) *Controller {
	return &Controller{
		serviceProvider: serviceProvider,
	}
}

func (ctrl *Controller) Register(c echo.Context) error {
	// Implement the registration logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the registration

	// Example response
	return c.JSON(200, map[string]string{"message": "Registration successful"})
}

func (ctrl *Controller) Login(c echo.Context) error {
	// Implement the login logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the login

	// Example response
	return c.JSON(200, map[string]string{"message": "Login successful"})
}

func (ctrl *Controller) CreateAccount(c echo.Context) error {
	// Implement the account creation logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the account creation

	// Example response
	return c.JSON(200, map[string]string{"message": "Account created successfully"})
}

func (ctrl *Controller) CreateCard(c echo.Context) error {
	// Implement the card creation logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the card creation

	// Example response
	return c.JSON(200, map[string]string{"message": "Card created successfully"})
}

func (ctrl *Controller) Transfer(c echo.Context) error {
	// Implement the transfer logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the transfer

	// Example response
	return c.JSON(200, map[string]string{"message": "Transfer successful"})
}

func (ctrl *Controller) GetCreditSchedule(c echo.Context) error {
	// Implement the credit schedule retrieval logic here
	// For example, you can bind the request data to a struct and validate it
	// Then, call the appropriate service method to handle the schedule retrieval

	// Example response
	return c.JSON(200, map[string]string{"message": "Credit schedule retrieved successfully"})
}
