package sqlext

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
	"time"
)

type DB interface {
	// Get выполняет запрос к базе данных, извлекает одну запись и заполняет ее в объект dest.
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Select выполняет SQL-запрос для выбора нескольких строк и заполняет переданный контейнер dest результатами.
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Exec выполняет команду в базе данных, не возвращающую строки, например INSERT, UPDATE или DELETE.
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Query выполняет запрос к базе данных, возвращая результаты в виде *sql.Rows.
	// ctx используется для управления временем выполнения запроса.
	// query задает SQL-запрос в виде строки, а args передает аргументы для подстановки в запрос.
	// Возвращает результаты запроса или ошибку в случае неудачи.
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// NamedExec выполняет запрос к базе данных, используя именованные параметры.
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	// Rebind преобразует запрос из "?" в bindvar типа драйвера базы данных.
	Rebind(query string) string

	// BindNamed привязывает запрос, используя bindvar тип драйвера базы данных.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	// WithTx выполняет функцию AtomicFn в контексте транзакции.
	WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error
}

type config struct {
	maxOpenConns int
	maxIdleConns int

	connLifeTime time.Duration
	connIdleTime time.Duration

	tracingEnabled bool
}

type ConnOption func(*config) error

// WithMaxConns устанавливает максимальные значения для числа открытых и простаивающих соединений в настройках config.
// Возвращает ошибку, если лимит простаивающих соединений превышает лимит открытых соединений.
func WithMaxConns(idle, open int) ConnOption {
	return func(c *config) error {
		if idle > open {
			return fmt.Errorf("ожидаемое количество простаивающих соединений не может быть больше, чем открытых (%v, %v)", idle, open)
		}

		c.maxIdleConns, c.maxOpenConns = idle, open
		return nil
	}
}

// WithConnTime задает время жизни соединения и время его бездействия.
// Возвращает ошибку, если idle меньше life или life отрицательное.
func WithConnTime(life, idle time.Duration) ConnOption {
	return func(c *config) error {
		if idle < life {
			return fmt.Errorf(
				"ожидаемое время бездействия соединения не может быть меньше, чем время его жизни (%v, %v)",
				idle,
				life,
			)
		}
		if life < 0 {
			return fmt.Errorf(
				"соединение не может иметь отрицательное время жизни, получено: %v",
				life,
			)
		}

		c.connLifeTime, c.connIdleTime = life, idle
		return nil
	}
}

// WithTracingEnabled задаёт опцию для включения или отключения трассировки соединения.
func WithTracingEnabled(enabled bool) ConnOption {
	return func(c *config) error {
		c.tracingEnabled = enabled
		return nil
	}
}
