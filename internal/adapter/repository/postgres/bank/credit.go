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

func (c CreditRepository) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	return c.db.WithTx(ctx, fn, opts...)
}
