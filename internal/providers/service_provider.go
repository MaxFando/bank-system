package providers

import (
	"github.com/MaxFando/bank-system/config"
	"github.com/MaxFando/bank-system/internal/core/bank/service/bank"
	"github.com/MaxFando/bank-system/internal/core/bank/service/user"
	"log/slog"
)

type ServiceProvider struct {
	logger *slog.Logger
	cfg    *config.Config

	UserService    *user.Service
	AuthService    *user.AuthService
	AccountService *bank.AccountService
	CardService    *bank.CardService
	CreditService  *bank.CreditService
}

func NewServiceProvider(logger *slog.Logger, cfg *config.Config) *ServiceProvider {
	return &ServiceProvider{logger: logger, cfg: cfg}
}

func (p *ServiceProvider) RegisterDependency(provider *RepositoryProvider) {
	p.UserService = user.NewService(p.logger, provider.userRepository)
	p.AuthService = user.NewAuthService(p.logger, provider.userRepository)
	p.AccountService = bank.NewAccountService(p.logger, provider.accountRepository)
	p.CardService = bank.NewCardService(p.logger, p.cfg, p.AccountService, provider.cardRepository, provider.transactionRepository)
	p.CreditService = bank.NewCreditService(p.logger, provider.creditRepository, p.AccountService)
}
