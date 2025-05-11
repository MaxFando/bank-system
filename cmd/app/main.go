package main

import (
	"context"
	"github.com/MaxFando/bank-system/config"
	"github.com/MaxFando/bank-system/internal/app"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()

	application := app.NewApp(cfg)
	err := application.Init(ctx)
	if err != nil {
		application.Logger().Error("Ошибка при загрузке конфигурации", err)
		return
	}

	if err := application.Run(ctx); err != nil {
		application.Logger().Error("Приложение завершило работу с ошибкой", err)
	}

	cancel()
	application.Shutdown(ctx)
	application.Logger().Info("Приложение завершило работу")
}
