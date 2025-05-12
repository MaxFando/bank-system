package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
)

type CreditRepository struct {
	db sqlext.DB
}

// NewCreditRepository создает новый экземпляр CreditRepository с заданным объектом базы данных.
func NewCreditRepository(db sqlext.DB) *CreditRepository {
	return &CreditRepository{
		db: db,
	}
}

func (c CreditRepository) Save(ctx context.Context, credit *entity.Credit) (*entity.Credit, error) {
	query := `INSERT INTO main.credits (user_id, amount, interest_rate, term_in_months) values ($1, $2, $3, $4) RETURNING id`

	err := c.db.Get(ctx, credit, query, credit.UserID, credit.Amount, credit.InterestRate, credit.TermInMonths)
	if err != nil {
		return nil, fmt.Errorf("failed to save credit: %w", err)
	}

	return credit, nil
}

func (c CreditRepository) CreatePaymentSchedule(ctx context.Context, paymentSchedule *entity.PaymentSchedule) (*entity.PaymentSchedule, error) {
	query := `INSERT INTO main.payment_schedules (credit_id, payment_date, payment_amount, principal_amount, interest_amount)
    				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := c.db.Get(
		ctx,
		paymentSchedule,
		query,
		paymentSchedule.CreditID,
		paymentSchedule.PaymentDate,
		paymentSchedule.PaymentAmount,
		paymentSchedule.PrincipalAmount,
		paymentSchedule.InterestAmount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment schedule: %w", err)
	}

	return paymentSchedule, nil
}

func (c CreditRepository) GetCreditByID(ctx context.Context, creditID int32) (*entity.Credit, error) {
	query := `SELECT id, user_id, amount, interest_rate, term_in_months FROM main.credits WHERE id = $1`
	credit := new(entity.Credit)

	err := c.db.Get(ctx, credit, query, creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit by ID: %w", err)
	}

	return credit, nil
}

func (c CreditRepository) GetPaymentSchedule(ctx context.Context, creditID int32) ([]entity.PaymentSchedule, error) {
	query := `SELECT id, credit_id, payment_date, payment_amount, principal_amount, interest_amount
				FROM main.payment_schedules WHERE credit_id = $1`

	var paymentSchedules []entity.PaymentSchedule
	err := c.db.Select(ctx, &paymentSchedules, query, creditID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment schedule: %w", err)
	}

	return paymentSchedules, nil
}

func (c CreditRepository) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	return c.db.WithTx(ctx, fn, opts...)
}
