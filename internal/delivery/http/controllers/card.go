package controllers

import (
	"github.com/MaxFando/bank-system/internal/core/bank/service/bank"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type CardController struct {
	accountService *bank.AccountService
	cardService    *bank.CardService
}

func NewCardController(accountService *bank.AccountService, cardService *bank.CardService) *CardController {
	return &CardController{
		accountService: accountService,
		cardService:    cardService,
	}
}

func (ctrl *CardController) CreateCard(c echo.Context) error {
	type request struct {
		AccountID int32 `json:"account_id" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	userID := c.Get("user_id").(int32)
	account, err := ctrl.accountService.GetAccountByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve account"})
	}

	if account.ID != req.AccountID {
		return c.JSON(403, map[string]string{"error": "Unauthorized account access"})
	}

	card, err := ctrl.cardService.Create(c.Request().Context(), account)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Card creation failed"})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Card created successfully",
		"card":    card,
	})
}

func (ctrl *CardController) GetCardsByAccountID(c echo.Context) error {
	type request struct {
		AccountID int32 `json:"account_id" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	userID := c.Get("user_id").(int32)
	account, err := ctrl.accountService.GetAccountByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve account"})
	}

	if account.ID != req.AccountID {
		return c.JSON(403, map[string]string{"error": "Unauthorized account access"})
	}

	cards, err := ctrl.cardService.FindByAccountID(c.Request().Context(), req.AccountID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve cards"})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Cards retrieved successfully",
		"cards":   cards,
	})
}

func (ctrl *CardController) Transfer(c echo.Context) error {
	type request struct {
		CardID          int32           `json:"card_id" validate:"required"`
		RecipientCardID int32           `json:"recipient_card_id" validate:"required"`
		Amount          decimal.Decimal `json:"amount" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	userID := c.Get("user_id").(int32)
	account, err := ctrl.accountService.GetAccountByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve account"})
	}

	card, err := ctrl.cardService.FindByID(c.Request().Context(), req.CardID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve card"})
	}

	if card.AccountID != account.ID {
		return c.JSON(403, map[string]string{"error": "Unauthorized card access"})
	}

	err = ctrl.cardService.Transfer(c.Request().Context(), req.CardID, req.RecipientCardID, req.Amount)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Transfer failed"})
	}

	return c.JSON(200, map[string]string{"message": "Transfer successful"})
}
