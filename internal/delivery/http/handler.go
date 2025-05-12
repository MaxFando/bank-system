package http

import (
	"context"
	"github.com/MaxFando/bank-system/internal/delivery/http/controllers"
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

	userController := controllers.NewUserController(provider.UserService, provider.AuthService)
	echoMainServer.POST("/register", userController.Register)
	echoMainServer.POST("/login", userController.Login)

	accountController := controllers.NewAccountController(provider.AccountService)
	echoMainServer.POST("/accounts", accountController.CreateAccount, echo.WrapMiddleware(auth.AuthMiddleware))

	cardController := controllers.NewCardController(provider.AccountService, provider.CardService)
	echoMainServer.POST("/cards", cardController.CreateCard, echo.WrapMiddleware(auth.AuthMiddleware))
	echoMainServer.POST("/cards/transfer", cardController.Transfer, echo.WrapMiddleware(auth.AuthMiddleware))

	creditController := controllers.NewCreditController(provider.CreditService)
	echoMainServer.POST("/credits", creditController.Create, echo.WrapMiddleware(auth.AuthMiddleware))
	echoMainServer.GET("/credits/:credit_id/schedule", creditController.GetCreditSchedule, echo.WrapMiddleware(auth.AuthMiddleware))

	return echoMainServer
}
