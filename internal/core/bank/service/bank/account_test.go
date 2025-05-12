package bank

import (
	"context"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

type NullWriter struct{}

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func TestAccountService_GetAccountByID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name     string
		id       int32
		mockFunc func(*MockAccountRepository)
		expected *entity.Account
		err      error
	}{
		{
			name: "successful account retrieval",
			id:   1,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().FindByID(gomock.Any(), int32(1)).Return(&entity.Account{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expected: &entity.Account{
				ID:     1,
				UserID: 1,
			},
		},
		{
			name: "error on account retrieval",
			id:   1,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().FindByID(gomock.Any(), int32(1)).Return(nil, assert.AnError)
			},
			err: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			tc.mockFunc(repo)

			service := NewAccountService(logger, repo)
			got, err := service.GetAccountByID(context.TODO(), tc.id)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestAccountService_GetAccountByUserID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name     string
		userID   int32
		mockFunc func(*MockAccountRepository)
		expected *entity.Account
		err      error
	}{
		{
			name:   "successful account retrieval by user ID",
			userID: 1,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().GetAccountByUserID(gomock.Any(), int32(1)).Return(&entity.Account{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expected: &entity.Account{
				ID:     1,
				UserID: 1,
			},
		},
		{
			name:   "user ID not found",
			userID: 2,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().GetAccountByUserID(gomock.Any(), int32(2)).Return(nil, nil)
			},
			expected: nil,
			err:      nil,
		},
		{
			name:   "error retrieving account by user ID",
			userID: 3,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().GetAccountByUserID(gomock.Any(), int32(3)).Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			tc.mockFunc(repo)

			service := NewAccountService(logger, repo)
			got, err := service.GetAccountByUserID(context.TODO(), tc.userID)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestAccountService_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name           string
		userID         int32
		initialBalance decimal.Decimal
		accountType    entity.AccountType
		mockFunc       func(*MockAccountRepository)
		expected       *entity.Account
		err            error
	}{
		{
			name:           "successful account creation",
			userID:         1,
			initialBalance: decimal.NewFromInt(1000),
			accountType:    entity.SavingsAccount,
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
					Return(&entity.Account{
						UserID:        1,
						AccountNumber: "ACC123",
						Balance:       decimal.NewFromInt(1000),
						AccountType:   entity.SavingsAccount,
					}, nil)
			},
			expected: &entity.Account{
				UserID:        1,
				AccountNumber: "ACC123",
				Balance:       decimal.NewFromInt(1000),
				AccountType:   entity.SavingsAccount,
			},
			err: nil,
		},
		{
			name:           "error creating account with negative balance",
			userID:         1,
			initialBalance: decimal.NewFromInt(-100),
			accountType:    entity.SavingsAccount,
			mockFunc:       nil,
			expected:       nil,
			err:            entity.ErrDepositNegativeAmount,
		},
		{
			name:           "error saving account",
			userID:         2,
			initialBalance: decimal.NewFromInt(500),
			accountType:    entity.AccountType("current"),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(&entity.Account{})).
					Return(nil, assert.AnError)
			},
			expected: nil,
			err:      assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			if tc.mockFunc != nil {
				tc.mockFunc(repo)
			}

			service := NewAccountService(logger, repo)
			got, err := service.Create(context.TODO(), tc.userID, tc.initialBalance, tc.accountType)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestAccountService_Deposit(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name      string
		accountID int32
		amount    decimal.Decimal
		mockFunc  func(m *MockAccountRepository)
		expected  error
	}{
		{
			name:      "successful deposit",
			accountID: 1,
			amount:    decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: nil,
		},
		{
			name:      "account not found",
			accountID: 2,
			amount:    decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name:      "negative deposit amount",
			accountID: 3,
			amount:    decimal.NewFromInt(-100),
			mockFunc:  nil,
			expected:  entity.ErrDepositNegativeAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			if tc.mockFunc != nil {
				tc.mockFunc(repo)
			}

			service := NewAccountService(logger, repo)
			err := service.Deposit(context.TODO(), tc.accountID, tc.amount)

			if tc.expected != nil {
				assert.ErrorAs(t, err, &tc.expected)
				assert.Equal(t, tc.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_Withdraw(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name      string
		accountID int32
		amount    decimal.Decimal
		mockFunc  func(m *MockAccountRepository)
		expected  error
	}{
		{
			name:      "successful withdrawal",
			accountID: 1,
			amount:    decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: nil,
		},
		{
			name:      "account not found",
			accountID: 2,
			amount:    decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name:      "insufficient balance",
			accountID: 3,
			amount:    decimal.NewFromInt(1000),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: entity.ErrInsufficientFunds,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			if tc.mockFunc != nil {
				tc.mockFunc(repo)
			}

			service := NewAccountService(logger, repo)
			err := service.Withdraw(context.TODO(), tc.accountID, tc.amount)

			if tc.expected != nil {
				assert.ErrorAs(t, err, &tc.expected)
				assert.Equal(t, tc.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name          string
		fromAccountID int32
		toAccountID   int32
		amount        decimal.Decimal
		mockFunc      func(m *MockAccountRepository)
		expected      error
	}{
		{
			name:          "successful transfer",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: nil,
		},
		{
			name:          "source account not found",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name:          "target account not found",
			fromAccountID: 1,
			toAccountID:   3,
			amount:        decimal.NewFromInt(100),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name:          "insufficient balance for transfer",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        decimal.NewFromInt(1000),
			mockFunc: func(m *MockAccountRepository) {
				m.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			expected: entity.ErrInsufficientFunds,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockAccountRepository(ctrl)
			if tc.mockFunc != nil {
				tc.mockFunc(repo)
			}

			service := NewAccountService(logger, repo)
			err := service.Transfer(context.TODO(), tc.fromAccountID, tc.toAccountID, tc.amount)

			if tc.expected != nil {
				assert.ErrorAs(t, err, &tc.expected)
				assert.Equal(t, tc.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
