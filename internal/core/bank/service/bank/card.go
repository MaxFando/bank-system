//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mock_${GOFILE}.go -package=${GOPACKAGE}
package bank

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/MaxFando/bank-system/config"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/openpgp"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

// CardRepository определяет методы для работы с хранилищем данных карт.
// Save сохраняет новую карту или обновляет существующую.
// FindByID ищет карту по уникальному идентификатору.
// FindByAccountID возвращает список карт, привязанных к определенному аккаунту.
type CardRepository interface {
	Save(ctx context.Context, card *entity.Card) (*entity.Card, error)
	FindByID(ctx context.Context, id int32) (*entity.Card, error)
	FindByAccountID(ctx context.Context, accountID int32) ([]entity.Card, error)
}

// CardTransactionRepository предоставляет методы для работы с операциями по картам, такими как перевод, снятие и пополнение.
type CardTransactionRepository interface {
	Transfer(ctx context.Context, cardID int32, amount decimal.Decimal) (int32, error)
	Withdraw(ctx context.Context, cardID int32, amount decimal.Decimal) (int32, error)
	Deposit(ctx context.Context, cardID int32, amount decimal.Decimal) (int32, error)
	FindByID(ctx context.Context, id int32) (*entity.CardTransaction, error)

	WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error
}

// CardService предоставляет методы для работы с банковскими картами, включая создание, обновление и обработку транзакций.
type CardService struct {
	cardRepository            CardRepository
	cardTransactionRepository CardTransactionRepository

	accountService *AccountService

	publicKeyPath  string
	privateKeyPath string
	passphrase     string

	logger *slog.Logger
}

// NewCardService создает и возвращает новый экземпляр CardService с заданными параметрами.
func NewCardService(
	logger *slog.Logger,
	cfg *config.Config,
	accountService *AccountService,
	cardRepository CardRepository,
	cardTransactionRepository CardTransactionRepository,
) *CardService {
	return &CardService{
		accountService:            accountService,
		cardRepository:            cardRepository,
		cardTransactionRepository: cardTransactionRepository,
		publicKeyPath:             cfg.PublicKeyPath,
		privateKeyPath:            cfg.PrivateKeyPath,
		passphrase:                cfg.Passphrase,
		logger:                    logger,
	}
}

// Create создает новую карту для указанного аккаунта, сохраняет данные карты и возвращает созданную карту или ошибку.
func (s *CardService) Create(ctx context.Context, account *entity.Account) (*entity.Card, error) {
	card := &entity.Card{
		AccountID:      account.ID,
		CardNumber:     generateCardNumber(account.ID),
		ExpirationDate: time.Now().AddDate(10, 0, 0),
		CVV:            generateCVV(),
		Status:         entity.CardActive,
	}

	// Генерация HMAC для проверки целостности данных карты
	hmac := s.generateHMAC(fmt.Sprintf("%s:%s:%s", card.CardNumber, card.CVV, card.ExpirationDate.Format("2006-01-02")), account.AccountNumber.String())
	card.HMAC = hmac

	// Шифрование данных карты
	encryptedData, err := s.encryptCardDataPGP(card, s.publicKeyPath, nil)
	if err != nil {
		s.logger.Error("failed to encrypt card data", "error", err)
		return nil, fmt.Errorf("failed to encrypt card data: %w", err)
	}

	card.EncryptedData = encryptedData

	// Сохранение карты в хранилище
	savedCard, err := s.cardRepository.Save(ctx, card)
	if err != nil {
		s.logger.Error("failed to save card", "error", err)
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	s.logger.Info("card created successfully", "card_id", savedCard.ID, "account_id", account.ID)
	return savedCard, nil
}

// FindByID находит и возвращает карту по указанному идентификатору, либо ошибку, если карта не найдена или произошла ошибка.
func (s *CardService) FindByID(ctx context.Context, id int32) (*entity.Card, error) {
	card, err := s.cardRepository.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to find card by ID", "error", err)
		return nil, fmt.Errorf("failed to find card by ID: %w", err)
	}

	e, err2 := s.smth(err, card)
	if err2 != nil {
		return e, err2
	}

	s.logger.Info("card found successfully", "card_id", card.ID)
	return card, nil
}

// FindByAccountID ищет и возвращает список карт, связанных с указанным идентификатором аккаунта.
// Возвращает ошибку, если карта не найдена или произошла проблема при обработке данных.
func (s *CardService) FindByAccountID(ctx context.Context, accountID int32) ([]entity.Card, error) {
	cards, err := s.cardRepository.FindByAccountID(ctx, accountID)
	if err != nil {
		s.logger.Error("failed to find cards by account ID", "error", err)
		return nil, fmt.Errorf("failed to find cards by account ID: %w", err)
	}

	s.logger.Info("cards found successfully", "account_id", accountID, "cards_count", len(cards))

	for i := range cards {
		_, err := s.smth(err, &cards[i])
		if err != nil {
			s.logger.Error("failed to decrypt card data", "error", err)
			return nil, fmt.Errorf("failed to decrypt card data: %w", err)
		}
	}

	return cards, nil
}

// smth проверяет целостность данных карты, расшифровывает их и обновляет информацию объекта карты.
func (s *CardService) smth(err error, card *entity.Card) (*entity.Card, error) {
	// Расшифровка данных карты
	decryptCardData, err := s.decryptCardData(card.EncryptedData, s.privateKeyPath, s.passphrase)
	if err != nil {
		s.logger.Error("failed to decrypt card data", "error", err)
		return nil, fmt.Errorf("failed to decrypt card data: %w", err)
	}

	// Проверка HMAC для целостности данных карты
	expectedHMAC := s.generateHMAC(decryptCardData, fmt.Sprintf("%d", card.AccountID))
	if card.HMAC != expectedHMAC {
		s.logger.Error("HMAC verification failed", "expected_hmac", expectedHMAC, "actual_hmac", card.HMAC)
		return nil, fmt.Errorf("HMAC verification failed")
	}

	decryptCardDataParts := bytes.Split([]byte(decryptCardData), []byte(":"))
	if len(decryptCardDataParts) != 3 {
		s.logger.Error("invalid decrypted card data format")
		return nil, fmt.Errorf("invalid decrypted card data format")
	}

	card.CardNumber = string(decryptCardDataParts[0])
	card.CVV = string(decryptCardDataParts[1])
	card.ExpirationDate, err = time.Parse("2006-01-02", string(decryptCardDataParts[2]))
	return nil, nil
}

// Transfer выполняет перевод указанной суммы с одной карты на другую в рамках заданного контекста.
func (s *CardService) Transfer(ctx context.Context, fromCardID, toCardID int32, amount decimal.Decimal) error {
	err := s.cardTransactionRepository.WithTx(ctx, func(ctx context.Context) error {
		fromCard, err := s.cardRepository.FindByID(ctx, fromCardID)
		if err != nil {
			s.logger.Error("failed to find source card", "error", err)
			return fmt.Errorf("failed to find source card: %w", err)
		}

		toCard, err := s.cardRepository.FindByID(ctx, toCardID)
		if err != nil {
			s.logger.Error("failed to find target card", "error", err)
			return fmt.Errorf("failed to find target card: %w", err)
		}

		fromAccount, err := s.accountService.GetAccountByID(ctx, fromCard.AccountID)
		if err != nil {
			s.logger.Error("failed to find source account", "error", err)
			return fmt.Errorf("failed to find source account: %w", err)
		}

		toAccount, err := s.accountService.GetAccountByID(ctx, toCard.AccountID)
		if err != nil {
			s.logger.Error("failed to find target account", "error", err)
			return fmt.Errorf("failed to find target account: %w", err)
		}

		if err := fromAccount.Transfer(toAccount, amount); err != nil {
			s.logger.Error("failed to transfer amount", "error", err)
			return fmt.Errorf("failed to transfer amount: %w", err)
		}

		transactionID, err := s.cardTransactionRepository.Transfer(ctx, fromCardID, amount)
		if err != nil {
			s.logger.Error("failed to transfer money", "error", err)
			return fmt.Errorf("failed to transfer money: %w", err)
		}

		s.logger.Info("money transferred successfully", "transaction_id", transactionID, "card_id", fromCardID, "amount", amount)
		return nil
	})

	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("transaction completed successfully", "card_id", fromCardID, "amount", amount)
	return nil
}

// Withdraw выполняет снятие указанной суммы с карты, идентифицированной cardID, с учетом контекста выполнения.
func (s *CardService) Withdraw(ctx context.Context, cardID int32, amount decimal.Decimal) error {
	err := s.cardTransactionRepository.WithTx(ctx, func(ctx context.Context) error {
		fromCard, err := s.cardRepository.FindByID(ctx, cardID)
		if err != nil {
			s.logger.Error("failed to find source card", "error", err)
			return fmt.Errorf("failed to find source card: %w", err)
		}

		fromAccount, err := s.accountService.GetAccountByID(ctx, fromCard.AccountID)
		if err != nil {
			s.logger.Error("failed to find source account", "error", err)
			return fmt.Errorf("failed to find source account: %w", err)
		}

		err = s.accountService.Withdraw(ctx, fromAccount.ID, amount)
		if err != nil {
			s.logger.Error("failed to withdraw amount", "error", err)
			return fmt.Errorf("failed to withdraw amount: %w", err)
		}

		transactionID, err := s.cardTransactionRepository.Withdraw(ctx, cardID, amount)
		if err != nil {
			s.logger.Error("failed to withdraw money", "error", err)
			return fmt.Errorf("failed to withdraw money: %w", err)
		}

		s.logger.Info("money withdrawn successfully", "transaction_id", transactionID, "card_id", cardID, "amount", amount)
		return nil
	})
	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	s.logger.Info("transaction completed successfully", "card_id", cardID, "amount", amount)
	return nil
}

// Deposit пополняет баланс карты на указанную сумму.
// Возвращает ошибку, если операция завершилась неуспешно.
func (s *CardService) Deposit(ctx context.Context, cardID int32, amount decimal.Decimal) error {
	err := s.cardTransactionRepository.WithTx(ctx, func(ctx context.Context) error {
		card, err := s.cardRepository.FindByID(ctx, cardID)
		if err != nil {
			s.logger.Error("failed to find source card", "error", err)
			return fmt.Errorf("failed to find source card: %w", err)
		}

		account, err := s.accountService.GetAccountByID(ctx, card.AccountID)
		if err != nil {
			s.logger.Error("failed to find source account", "error", err)
			return fmt.Errorf("failed to find source account: %w", err)
		}

		err = s.accountService.Deposit(ctx, account.ID, amount)
		if err != nil {
			s.logger.Error("failed to deposit amount", "error", err)
			return fmt.Errorf("failed to deposit amount: %w", err)
		}

		transactionID, err := s.cardTransactionRepository.Deposit(ctx, cardID, amount)
		if err != nil {
			s.logger.Error("failed to deposit money", "error", err)
			return fmt.Errorf("failed to deposit money: %w", err)
		}

		s.logger.Info("money deposited successfully", "transaction_id", transactionID, "card_id", cardID, "amount", amount)
		return nil
	})
	if err != nil {
		s.logger.Error("transaction failed", "error", err)
		return fmt.Errorf("transaction failed: %w", err)
	}
	s.logger.Info("transaction completed successfully", "card_id", cardID, "amount", amount)
	return nil
}

// generateCVV генерирует случайный CVV код длиной 3 символа в виде строки.
func generateCVV() string {
	rand.Seed(time.Now().UnixNano())

	cvvLength := 3

	var cvv int
	cvv = rand.Intn(900) + 100

	return fmt.Sprintf("%0*d", cvvLength, cvv)
}

// generateCardNumber создает и возвращает уникальный номер карты на основе ID аккаунта.
// Для генерации используется алгоритм Luhn для проверки корректности номера.
func generateCardNumber(accountID int32) string {
	seed := int64(accountID) + time.Now().UnixNano()
	rand.Seed(seed)

	cardNumber := make([]int, 16)
	for i := 0; i < 15; i++ {
		cardNumber[i] = rand.Intn(10)
	}

	checksum := 0
	for i := 0; i < 15; i++ {
		digit := cardNumber[i]
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		checksum += digit
	}

	cardNumber[15] = (10 - (checksum % 10)) % 10

	cardNumberStr := ""
	for _, digit := range cardNumber {
		cardNumberStr += fmt.Sprintf("%d", digit)
	}

	return cardNumberStr
}

// generateCardNumber генерирует уникальный номер карты с использованием алгоритма Луна для проверки контрольной цифры.
func (s *CardService) generateCardNumber() string {
	cardNumber := make([]int, 16)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 15; i++ {
		cardNumber[i] = rand.Intn(10)
	}

	checksum := 0
	for i := 0; i < 15; i++ {
		digit := cardNumber[i]
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		checksum += digit
	}

	cardNumber[15] = (10 - (checksum % 10)) % 10

	cardNumberStr := ""
	for _, digit := range cardNumber {
		cardNumberStr += fmt.Sprintf("%d", digit)
	}

	return cardNumberStr
}

// encryptCardDataPGP шифрует данные карты, используя PGP и указанный публичный ключ. Возвращает зашифрованную строку или ошибку.
func (s *CardService) encryptCardDataPGP(card *entity.Card, publicKeyPath string, signedEntity *openpgp.Entity) (string, error) {
	pubKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not open public key file: %v", err)
	}
	defer pubKeyFile.Close()

	pubKeyRing, err := openpgp.ReadArmoredKeyRing(pubKeyFile)
	if err != nil {
		return "", fmt.Errorf("could not read public key: %v", err)
	}

	cardData := fmt.Sprintf("%s:%s:%s", card.CardNumber, card.CVV, card.ExpirationDate)

	var encryptedData bytes.Buffer

	plaintext, err := openpgp.Encrypt(&encryptedData, pubKeyRing, signedEntity, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not start encryption: %v", err)
	}

	_, err = io.WriteString(plaintext, cardData)
	if err != nil {
		return "", fmt.Errorf("could not write card data to encrypted stream: %v", err)
	}

	err = plaintext.Close()
	if err != nil {
		return "", fmt.Errorf("could not close encryption stream: %v", err)
	}

	return encryptedData.String(), nil
}

// decryptCardData расшифровывает данные карты, зашифрованные с использованием PGP, с использованием указанного приватного ключа.
// Принимает зашифрованные данные, путь к файлу приватного ключа и парольную фразу для дешифровки.
// Возвращает строку с расшифрованными данными или ошибку, если проблема возникла во время процесса.
func (s *CardService) decryptCardData(encryptedData string, privateKeyPath string, passphrase string) (string, error) {
	privKeyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not open private key file: %v", err)
	}
	defer privKeyFile.Close()

	privKeyRing, err := openpgp.ReadArmoredKeyRing(privKeyFile)
	if err != nil {
		return "", fmt.Errorf("could not read private key: %v", err)
	}

	encryptedReader := bytes.NewReader([]byte(encryptedData))

	md, err := openpgp.ReadMessage(encryptedReader, privKeyRing, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not read encrypted message: %v", err)
	}

	if passphrase != "" {
		err = privKeyRing[0].PrivateKey.Decrypt([]byte(passphrase))
		if err != nil {
			return "", fmt.Errorf("could not decrypt private key: %v", err)
		}
	}

	var decryptedData bytes.Buffer
	_, err = io.Copy(&decryptedData, md.UnverifiedBody)
	if err != nil {
		return "", fmt.Errorf("could not copy decrypted data: %v", err)
	}

	return decryptedData.String(), nil
}

// generateHMAC генерирует HMAC-хэш строки `data` с использованием заданного секретного ключа `secretKey`.
func (s *CardService) generateHMAC(data string, secretKey string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
