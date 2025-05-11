//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mock_${GOFILE}.go -package=${GOPACKAGE}
package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/shopspring/decimal"
	"log/slog"
	"math/rand"
	"time"
)

// AccountRepository задает интерфейс для операций с банковскими счетами, включая сохранение и поиск по идентификатору.
type AccountRepository interface {
	Save(ctx context.Context, account *entity.Account) (*entity.Account, error)
	FindByID(ctx context.Context, id int32) (*entity.Account, error)
}

// AccountService предоставляет методы для работы с банковскими счетами, включая создание, пополнение, снятие и переводы.
type AccountService struct {
	repo   AccountRepository
	logger *slog.Logger
}

// NewAccountService создает новый экземпляр AccountService с указанным логгером и репозиторием.
func NewAccountService(logger *slog.Logger, repo AccountRepository) *AccountService {
	return &AccountService{
		repo:   repo,
		logger: logger,
	}
}

// Create создает новый банковский счет с указанными параметрами и сохраняет его в хранилище.
func (s *AccountService) Create(
	ctx context.Context,
	userID int32,
	initialBalance decimal.Decimal,
	accountType entity.AccountType,
) (*entity.Account, error) {
	if initialBalance.LessThan(decimal.Zero) {
		return nil, entity.ErrDepositNegativeAmount
	}

	number, err := generateAccountNumber(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate account number: %w", err)
	}

	account := &entity.Account{
		UserID:        userID,
		AccountNumber: number,
		Balance:       initialBalance,
		Currency:      entity.RUB,
		AccountType:   accountType,
	}

	createdAccount, err := s.repo.Save(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	s.logger.Info("Account created successfully", "account_number", createdAccount.AccountNumber)
	return createdAccount, nil
}

// generateAccountNumber генерирует уникальный номер счета на основе userID и случайного числа.
// Возвращает номер счета и ошибку, если номер не прошел проверку валидности.
func generateAccountNumber(userID int32) (entity.AccountNumber, error) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	randomPart := randGen.Intn(1000000000000000)
	randomPartStr := fmt.Sprintf("%015d", randomPart)

	accountNumber := entity.AccountNumber(fmt.Sprintf("%d%s", userID, randomPartStr))

	if err := accountNumber.Validate(); err != nil {
		return "", err
	}

	return accountNumber, nil
}

// Deposit выполняет пополнение баланса указанного счета на заданную сумму. Возвращает ошибку, если операция завершилась неудачей.
func (s *AccountService) Deposit(ctx context.Context, accountID int32, amount decimal.Decimal) error {
	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	if err := account.Deposit(amount); err != nil {
		return fmt.Errorf("failed to deposit amount: %w", err)
	}

	if _, err := s.repo.Save(ctx, account); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	s.logger.Info("Deposit successful", "account_number", account.AccountNumber, "amount", amount)
	return nil
}

// Withdraw выполняет операцию снятия указанной суммы со счета. Возвращает ошибку, если операция невозможна.
func (s *AccountService) Withdraw(ctx context.Context, accountID int32, amount decimal.Decimal) error {
	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	if err := account.Withdraw(amount); err != nil {
		return fmt.Errorf("failed to withdraw amount: %w", err)
	}

	if _, err := s.repo.Save(ctx, account); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	s.logger.Info("Withdrawal successful", "account_number", account.AccountNumber, "amount", amount)
	return nil
}

// Transfer выполняет перевод суммы между указанными счетами. Возвращает ошибку в случае неудачи.
func (s *AccountService) Transfer(ctx context.Context, fromAccountID, toAccountID int32, amount decimal.Decimal) error {
	fromAccount, err := s.repo.FindByID(ctx, fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to find source account: %w", err)
	}

	toAccount, err := s.repo.FindByID(ctx, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to find target account: %w", err)
	}

	if err := fromAccount.Transfer(toAccount, amount); err != nil {
		return fmt.Errorf("failed to transfer amount: %w", err)
	}

	if _, err := s.repo.Save(ctx, fromAccount); err != nil {
		return fmt.Errorf("failed to update source account balance: %w", err)
	}

	if _, err := s.repo.Save(ctx, toAccount); err != nil {
		return fmt.Errorf("failed to update target account balance: %w", err)
	}

	s.logger.Info("Transfer successful", "from_account_number", fromAccount.AccountNumber, "to_account_number", toAccount.AccountNumber, "amount", amount)
	return nil
}
