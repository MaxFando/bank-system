package controllers

import (
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/internal/core/bank/service/bank"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type AccountController struct {
	accountService *bank.AccountService
}

func NewAccountController(accountService *bank.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (ctrl *AccountController) CreateAccount(c echo.Context) error {
	type request struct {
		InitialBalance decimal.Decimal `json:"initial_balance" validate:"required,gte=0"`
		AccountType    string          `json:"account_type" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	userID := c.Get("user_id").(int32)
	account, err := ctrl.accountService.Create(c.Request().Context(), userID, req.InitialBalance, entity.AccountType(req.AccountType))
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Account creation failed"})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Account created successfully",
		"account": account,
	})
}
