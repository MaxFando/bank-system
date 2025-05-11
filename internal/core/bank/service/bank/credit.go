package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
	"github.com/shopspring/decimal"
	"log/slog"
	"time"
)

// CreditRepository предоставляет методы для работы с хранилищем кредитов.
// Включает операции сохранения, получения и создания графиков платежей.
// Также поддерживает выполнение операций в транзакциях.
type CreditRepository interface {
	Save(ctx context.Context, credit *entity.Credit) (*entity.Credit, error)
	CreatePaymentSchedule(ctx context.Context, paymentSchedule *entity.PaymentSchedule) (*entity.PaymentSchedule, error)
	GetCreditByID(ctx context.Context, creditID int32) (*entity.Credit, error)

	WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error
}

// CreditService предоставляет функциональность для управления кредитами, включая создание, расчет и хранение данных.
type CreditService struct {
	creditRepository CreditRepository
	accountService   *AccountService
	logger           *slog.Logger
}

// NewCreditService создает новый экземпляр CreditService с заданным логгером, репозиторием кредита и сервисом счетов.
func NewCreditService(logger *slog.Logger, creditRepository CreditRepository, accountService *AccountService) *CreditService {
	return &CreditService{
		creditRepository: creditRepository,
		accountService:   accountService,
		logger:           logger,
	}
}

// Create создает новый кредит для указанного пользователя и составляет график платежей на основе переданных параметров.
func (s *CreditService) Create(ctx context.Context, userID int32, principal, interestRate decimal.Decimal, termMonths int32) (*entity.Credit, error) {
	annuityPayment, err := calculateAnnuityPayment(principal, interestRate, int(termMonths))
	if err != nil {
		s.logger.Error("failed to calculate annuity payment", "error", err)
		return nil, fmt.Errorf("failed to calculate annuity payment: %w", err)
	}

	// Создание кредита
	credit := new(entity.Credit)
	err = s.creditRepository.WithTx(ctx, func(ctx context.Context) error {
		credit = &entity.Credit{
			UserID:       userID,
			Amount:       principal,
			InterestRate: interestRate,
			TermInMonths: termMonths,
		}
		createdCredit, err := s.creditRepository.Save(ctx, credit)
		if err != nil {
			s.logger.Error("failed to create credit", "error", err)
			return err
		}

		s.logger.Info("credit created successfully", "credit_id", createdCredit.ID)
		credit = createdCredit

		// Создание графика платежей
		for i := 0; i < int(termMonths); i++ {
			paymentSchedule := &entity.PaymentSchedule{
				CreditID:        createdCredit.ID,
				PaymentDate:     time.Now().AddDate(0, 0, (i+1)*30),
				PaymentAmount:   annuityPayment,
				PrincipalAmount: principal,
				InterestAmount:  interestRate,
				Balance:         principal.Sub(annuityPayment.Mul(decimal.NewFromInt(int64(i)))),
			}

			_, err = s.creditRepository.CreatePaymentSchedule(ctx, paymentSchedule)
			if err != nil {
				s.logger.Error("failed to create payment schedule", "error", err)
				return fmt.Errorf("failed to create payment schedule: %w", err)
			}
		}

		s.logger.Info("payment schedule created successfully", "credit_id", credit.ID)
		return nil
	})
	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("credit created successfully", "credit_id", credit.ID)
	return credit, nil
}

// WithdrawPayment выполняет списание указанной суммы `amount` из кредитного счёта `credit` с проверкой доступности средств.
// Возвращает ошибку, если операция не может быть завершена, например, из-за отсутствия средств на счету.
// Использует транзакцию для обеспечения согласованности данных между кредитным счётом и пользователем.
func (s *CreditService) WithdrawPayment(ctx context.Context, credit *entity.Credit, amount decimal.Decimal) error {
	err := s.creditRepository.WithTx(ctx, func(ctx context.Context) error {
		account, err := s.accountService.GetAccountByUserID(ctx, credit.UserID)
		if err != nil {
			s.logger.Error("failed to get account", "error", err)
			return fmt.Errorf("failed to get account: %w", err)
		}

		err = s.accountService.Withdraw(ctx, account.ID, amount)
		if err != nil {
			s.logger.Error("failed to withdraw amount", "error", err)
			return fmt.Errorf("failed to withdraw amount: %w", err)
		}
		s.logger.Info("withdrawal successful", "account_id", account.ID, "amount", amount)

		err = credit.Withdraw(amount)
		if err != nil {
			s.logger.Error("failed to withdraw from credit", "error", err)
			return fmt.Errorf("failed to withdraw from credit: %w", err)
		}

		s.logger.Info("withdrawal from credit successful", "credit_id", credit.ID, "amount", amount)

		// Обновление кредита
		_, err = s.creditRepository.Save(ctx, credit)
		if err != nil {
			s.logger.Error("failed to update credit", "error", err)
			return fmt.Errorf("failed to update credit: %w", err)
		}

		s.logger.Info("credit updated successfully", "credit_id", credit.ID)
		return nil
	})

	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("payment withdrawn successfully", "credit_id", credit.ID, "amount", amount)
	return nil
}

// ApplyPenalty накладывает штраф на кредит, вычитая 10% от суммы кредита с банковского счета пользователя.
func (s *CreditService) ApplyPenalty(ctx context.Context, credit *entity.Credit) error {
	penaltyAmount := credit.Amount.Mul(decimal.NewFromFloat(0.1)) // 10% penalty
	err := s.creditRepository.WithTx(ctx, func(ctx context.Context) error {
		account, err := s.accountService.GetAccountByUserID(ctx, credit.UserID)
		if err != nil {
			s.logger.Error("failed to get account", "error", err)
			return fmt.Errorf("failed to get account: %w", err)
		}

		err = s.accountService.Withdraw(ctx, account.ID, penaltyAmount)
		if err != nil {
			s.logger.Error("failed to withdraw penalty amount", "error", err)
			return fmt.Errorf("failed to withdraw penalty amount: %w", err)
		}
		s.logger.Info("penalty withdrawal successful", "account_id", account.ID, "amount", penaltyAmount)

		return nil
	})

	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("penalty applied successfully", "credit_id", credit.ID, "penalty_amount", penaltyAmount)
	return nil
}

// calculateAnnuityPayment вычисляет аннуитетный платеж для кредита на основе суммы, процентной ставки и срока в месяцах.
// Возвращает рассчитанный аннуитетный платеж или ошибку в случае некорректных входных данных.
func calculateAnnuityPayment(principal, interestRate decimal.Decimal, termMonths int) (decimal.Decimal, error) {
	termMonthsDecimal := decimal.NewFromInt(int64(termMonths))

	monthlyInterestRate := interestRate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))

	// Расчет аннуитетного платежа
	// A = P * (r * (1 + r)^n) / ((1 + r)^n - 1)
	onePlusR := decimal.NewFromInt(1).Add(monthlyInterestRate)
	onePlusRToN := onePlusR.Pow(termMonthsDecimal)
	numerator := monthlyInterestRate.Mul(onePlusRToN)
	denominator := onePlusRToN.Sub(decimal.NewFromInt(1))

	annuityPayment := principal.Mul(numerator).Div(denominator)

	return annuityPayment, nil
}
