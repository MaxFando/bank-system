package app

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/config"
	"github.com/MaxFando/bank-system/internal/delivery/http"
	"github.com/MaxFando/bank-system/internal/providers"
	"github.com/MaxFando/bank-system/pkg/sqlext"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger          *slog.Logger
	config          *config.Config
	serviceProvider *providers.ServiceProvider
	httpServer      *http.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (a *App) Logger() *slog.Logger {
	return a.logger
}

func (a *App) Init(ctx context.Context) error {
	if err := a.initProviders(ctx); err != nil {
		return fmt.Errorf("failed to init providers: %w", err)
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting application...")

	httpServer := http.NewHttpServer(http.NewHandler(ctx, a.serviceProvider), ":"+a.config.Port)
	httpServer.Serve()

	a.httpServer = httpServer
	a.logger.Info("HTTP server started", "port", a.config.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		return fmt.Errorf("signal: %s", s.String())
	case err := <-httpServer.Notify():
		return fmt.Errorf("server: %w", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	_ = a.httpServer.Shutdown()
}

func (a *App) initProviders(ctx context.Context) error {
	postgresConn, err := sqlext.NewPostgresDB(ctx, a.config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	repositoryProvider := providers.NewRepositoryProvider(postgresConn)
	repositoryProvider.RegisterDependency()

	serviceProvider := providers.NewServiceProvider(a.logger, a.config)
	serviceProvider.RegisterDependency(repositoryProvider)

	a.serviceProvider = serviceProvider

	return nil
}
