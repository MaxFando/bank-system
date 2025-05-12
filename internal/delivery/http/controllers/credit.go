package controllers

import (
	"github.com/MaxFando/bank-system/internal/core/bank/service/bank"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type CreditController struct {
	creditService *bank.CreditService
}

func NewCreditController(creditService *bank.CreditService) *CreditController {
	return &CreditController{
		creditService: creditService,
	}
}

func (ctrl *CreditController) Create(c echo.Context) error {
	type request struct {
		Principal decimal.Decimal `json:"principal" validate:"required"`
		Term      int32           `json:"term" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	userID := c.Get("user_id").(int32)
	credit, err := ctrl.creditService.Create(c.Request().Context(), userID, req.Principal, req.Term)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Credit creation failed"})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Credit created successfully",
		"credit":  credit,
	})
}

func (ctrl *CreditController) GetCreditSchedule(c echo.Context) error {
	type request struct {
		CreditID int32 `param:"credit_id" validate:"required"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, map[string]string{"error": "Validation failed"})
	}

	schedule, err := ctrl.creditService.GetPaymentSchedule(c.Request().Context(), req.CreditID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to get credit schedule"})
	}

	return c.JSON(200, schedule)
}
