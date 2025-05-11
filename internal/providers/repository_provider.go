package providers

import (
	"github.com/MaxFando/bank-system/internal/adapter/repository/postgres/bank"
	"github.com/MaxFando/bank-system/internal/adapter/repository/postgres/user"
	"github.com/MaxFando/bank-system/pkg/sqlext"
)

type RepositoryProvider struct {
	db sqlext.DB

	userRepository        *user.Repository
	accountRepository     *bank.AccountRepository
	cardRepository        *bank.CardRepository
	creditRepository      *bank.CreditRepository
	transactionRepository *bank.CardTransactionRepository
}

func NewRepositoryProvider(db sqlext.DB) *RepositoryProvider {
	return &RepositoryProvider{
		db: db,
	}
}

func (p *RepositoryProvider) RegisterDependency() {
	p.userRepository = user.NewRepository(p.db)
	p.accountRepository = bank.NewAccountRepository(p.db)
	p.cardRepository = bank.NewCardRepository(p.db)
	p.creditRepository = bank.NewCreditRepository(p.db)
	p.transactionRepository = bank.NewCardTransactionRepository(p.db)
}
