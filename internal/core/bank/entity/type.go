package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID          int32  `db:"id"`
	Email       string `db:"email"`
	Password    string `db:"password_hash"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	DateOfBirth string `db:"date_of_birth"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type AccountType string

const (
	SavingsAccount  AccountType = "savings"
	CheckingAccount AccountType = "checking"
	CreditAccount   AccountType = "credit"
)

type Account struct {
	ID            int32           `db:"id"`
	UserID        int32           `db:"user_id"`
	AccountNumber AccountNumber   `db:"account_number"`
	Balance       decimal.Decimal `db:"balance"`
	Currency      Currency        `db:"currency"`
	AccountType   AccountType     `db:"account_type"`
	CreatedAt     string          `db:"created_at"`
	UpdatedAt     string          `db:"updated_at"`
}

func (a *Account) Transfer(target *Account, amount decimal.Decimal) error {
	if err := a.Withdraw(amount); err != nil {
		return err
	}

	return target.Deposit(amount)
}

func (a *Account) Withdraw(amount decimal.Decimal) error {
	if a.Balance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	a.Balance = a.Balance.Sub(amount)

	return nil
}

func (a *Account) Deposit(amount decimal.Decimal) error {
	if amount.LessThan(decimal.Zero) {
		return ErrDepositNegativeAmount
	}

	a.Balance = a.Balance.Add(amount)

	return nil
}

type CardStatus string

const (
	CardActive   CardStatus = "active"
	CardInactive CardStatus = "inactive"
	CardBlocked  CardStatus = "blocked"
	CardExpired  CardStatus = "expired"
)

type Card struct {
	ID             int32 `db:"id"`
	AccountID      int32 `db:"account_id"`
	CardNumber     string
	ExpirationDate time.Time
	CVV            string
	Status         CardStatus
	EncryptedData  string `db:"encrypted_data"`
	HMAC           string `db:"hash"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

type TransactionType string

const (
	PaymentTransaction    TransactionType = "payment"
	WithdrawalTransaction TransactionType = "withdrawal"
	DepositTransaction    TransactionType = "deposit"
	TransferTransaction   TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionSuccess TransactionStatus = "success"
	TransactionFailed  TransactionStatus = "failed"
	TransactionPending TransactionStatus = "pending"
)

type CardTransaction struct {
	ID              int32             `db:"id"`
	CardID          int32             `db:"card_id"`
	Amount          decimal.Decimal   `db:"amount"`
	TransactionType TransactionType   `db:"transaction_type"`
	TransactionDate string            `db:"transaction_date"`
	Status          TransactionStatus `db:"status"`
}

type Transfer struct {
	ID            int32             `db:"id"`              // Идентификатор перевода
	FromAccountID int32             `db:"from_account_id"` // Внешний ключ на отправляющий счет
	ToAccountID   int32             `db:"to_account_id"`   // Внешний ключ на получающий счет
	Amount        decimal.Decimal   `db:"amount"`          // Сумма перевода
	TransferDate  string            `db:"transfer_date"`   // Дата перевода
	Status        TransactionStatus `db:"status"`          // Статус перевода (успешно, отклонено)
}

type CreditStatus string

const (
	CreditStatusActive CreditStatus = "active"
	CreditStatusPaid   CreditStatus = "paid"
)

type Credit struct {
	ID           int32           `db:"id"`             // Идентификатор кредита
	UserID       int32           `db:"user_id"`        // Внешний ключ на пользователя
	Amount       decimal.Decimal `db:"amount"`         // Сумма кредита
	InterestRate decimal.Decimal `db:"interest_rate"`  // Процентная ставка
	TermInMonths int32           `db:"term_in_months"` // Срок кредита (в месяцах)
	Status       CreditStatus    `db:"status"`         // Статус кредита (оформлен, погашен)
	CreatedAt    string          `db:"created_at"`     // Дата оформления кредита
	UpdatedAt    string          `db:"updated_at"`     // Дата последнего обновления
}

func (c *Credit) Withdraw(amount decimal.Decimal) error {
	if c.Status != CreditStatusActive {
		return ErrCreditNotActive
	}

	if amount.GreaterThan(c.Amount) {
		return ErrCreditAmountExceeded
	}

	c.Amount = c.Amount.Sub(amount)

	return nil
}

type PaymentSchedule struct {
	ID              int32           `db:"id"`               // Идентификатор записи
	CreditID        int32           `db:"credit_id"`        // Внешний ключ на кредит
	PaymentDate     time.Time       `db:"payment_date"`     // Дата платежа
	PaymentAmount   decimal.Decimal `db:"payment_amount"`   // Сумма платежа
	PrincipalAmount decimal.Decimal `db:"principal_amount"` // Сумма, погашенная по телу кредита
	InterestAmount  decimal.Decimal `db:"interest_amount"`  // Сумма, погашенная по процентам
	Penalty         decimal.Decimal `db:"penalty"`          // Штраф за просрочку
	Balance         decimal.Decimal `db:"balance"`          // Остаток долга после платежа
	CreatedAt       string          `db:"created_at"`       // Дата создания записи
	UpdatedAt       string          `db:"updated_at"`       // Дата последнего обновления
}

type FinancialTransaction struct {
	ID                int32             `db:"id"`                 // Идентификатор операции
	UserID            int32             `db:"user_id"`            // Внешний ключ на пользователя
	TransactionType   TransactionType   `db:"transaction_type"`   // Тип операции (пополнение, снятие и т.д.)
	Amount            decimal.Decimal   `db:"amount"`             // Сумма операции
	TransactionDate   time.Time         `db:"transaction_date"`   // Дата операции
	TransactionStatus TransactionStatus `db:"transaction_status"` // Статус операции (успешно, отклонено)
	CreatedAt         time.Time         `db:"created_at"`         // Дата создания записи
	UpdatedAt         time.Time         `db:"updated_at"`         // Дата последнего обновления
}

type NotificationType string

const (
	EmailNotification NotificationType = "email"
)

type NotificationStatus string

const (
	NotificationSent      NotificationStatus = "sent"
	NotificationNotSent   NotificationStatus = "not_sent"
	NotificationPending   NotificationStatus = "pending"
	NotificationFailed    NotificationStatus = "failed"
	NotificationDelivered NotificationStatus = "delivered"
	NotificationRead      NotificationStatus = "read"
)

type Notification struct {
	ID        int32              `db:"id"`         // Идентификатор уведомления
	UserID    int32              `db:"user_id"`    // Внешний ключ на пользователя
	Type      NotificationType   `db:"type"`       // Тип уведомления (email, SMS и т.д.)
	Subject   string             `db:"subject"`    // Тема письма
	Body      string             `db:"body"`       // Текст письма
	Status    NotificationStatus `db:"status"`     // Статус отправки (отправлено, не отправлено)
	CreatedAt time.Time          `db:"created_at"` // Дата отправки
}

type CentralBankRate struct {
	ID        int32           `db:"id"`         // Идентификатор записи
	Rate      decimal.Decimal `db:"rate"`       // Ключевая ставка
	RateDate  time.Time       `db:"rate_date"`  // Дата ставки
	CreatedAt time.Time       `db:"created_at"` // Дата записи
}
