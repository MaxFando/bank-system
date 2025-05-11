package cb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"io"
	"net/http"
	"time"
)

const (
	CentralBankServiceURL = "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"
	ContentTypeHeader     = "application/soap+xml; charset=utf-8"
	SOAPActionHeader      = "http://web.cbr.ru/KeyRate"
	MarginValue           = 5.0
)

// GetCentralBankRateWithMargin возвращает текущую ключевую ставку центрального банка с добавленной маржей.
func GetCentralBankRateWithMargin() (float64, error) {
	soapRequest := buildSOAPRequest()
	rawBody, err := sendRequest(soapRequest)
	if err != nil {
		return 0, fmt.Errorf("ошибка при отправке SOAP-запроса: %w", err)
	}
	rate, err := parseXMLResponse(rawBody)
	if err != nil {
		return 0, fmt.Errorf("ошибка при обработке XML-ответа: %w", err)
	}
	// Добавление маржи
	return rate + MarginValue, nil
}

// sendRequest отправляет SOAP-запрос к заданному сервису и возвращает необработанный ответ или ошибку.
func sendRequest(soapRequest string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", CentralBankServiceURL, bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}
	// Установка заголовков
	req.Header.Set("Content-Type", ContentTypeHeader)
	req.Header.Set("SOAPAction", SOAPActionHeader)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()
	// Чтение ответа
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}
	return rawBody, nil
}

// buildSOAPRequest создает строку SOAP-запроса для получения ключевой ставки из Центрального Банка России.
func buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
            <soap12:Body>
                <KeyRate xmlns="http://web.cbr.ru/">
                    <fromDate>%s</fromDate>
                    <ToDate>%s</ToDate>
                </KeyRate>
            </soap12:Body>
        </soap12:Envelope>`, fromDate, toDate)
}

// parseXMLResponse парсит XML-ответ и извлекает значение ключевой ставки из тега Rate.
func parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("ошибка парсинга XML: %w", err)
	}
	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, errors.New("данные ставки не найдены в ответе")
	}
	rateElement := krElements[0].FindElement("./Rate")
	if rateElement == nil {
		return 0, errors.New("тег Rate отсутствует в ответе")
	}
	var rate float64
	if _, err := fmt.Sscanf(rateElement.Text(), "%f", &rate); err != nil {
		return 0, fmt.Errorf("ошибка конвертации ставки в число: %w", err)
	}
	return rate, nil
}
