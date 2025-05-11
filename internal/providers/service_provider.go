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

	userService    *user.Service
	authService    *user.AuthService
	accountService *bank.AccountService
	cardService    *bank.CardService
	creditService  *bank.CreditService
}

func NewServiceProvider(logger *slog.Logger, cfg *config.Config) *ServiceProvider {
	return &ServiceProvider{logger: logger, cfg: cfg}
}

func (p *ServiceProvider) RegisterDependency(provider *RepositoryProvider) {
	p.userService = user.NewService(p.logger, provider.userRepository)
	p.authService = user.NewAuthService(p.logger, provider.userRepository)
	p.accountService = bank.NewAccountService(p.logger, provider.accountRepository)
	p.cardService = bank.NewCardService(p.logger, p.cfg, provider.cardRepository, provider.transactionRepository)
	p.creditService = bank.NewCreditService(p.logger, provider.creditRepository, p.accountService)
}
