package repository

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/1abobik1/zadanie_3/internal/entities"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
)

type CurrencyRepository interface {
	GetRates(date string) ([]entities.CurrencyRate, error)
}

type CBRApiClient struct {
	baseURL string
}

func NewCBRApiClient() *CBRApiClient {
	return &CBRApiClient{
		baseURL: "https://www.cbr.ru/scripts/XML_daily_eng.asp",
	}
}

func (c *CBRApiClient) GetRates(date string) ([]entities.CurrencyRate, error) {
	logFields := logrus.Fields{
		"date":   date,
		"method": "GetRates",
		"layer":  "repository",
	}

	url := fmt.Sprintf("%s?date_req=%s", c.baseURL, date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// Повторно используйте http.Client для keep-alive
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logFields).WithError(err).Error("Ошибка запроса к API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logFields).Warnf("Unexpected API response: %d %s",
			resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logFields).WithError(err).Error("Failed to read response body")
		return nil, err
	}

	var result struct {
		Date  string `xml:"Date,attr"`
		Rates []struct {
			CharCode string `xml:"CharCode"`
			Name     string `xml:"Name"`
			Value    string `xml:"Value"`
		} `xml:"Valute"`
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
	if err := decoder.Decode(&result); err != nil {
		logrus.WithFields(logFields).WithError(err).Error("XML parsing failed")
		return nil, err
	}

	parsedDate, err := time.Parse("02.01.2006", result.Date)
	if err != nil {
		logrus.WithFields(logFields).WithError(err).Error("Date parsing failed")
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}
	dateStr := parsedDate.Format("2006-01-02")

	var rates []entities.CurrencyRate
	for _, r := range result.Rates {
		valueStr := strings.Replace(r.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, err
		}
		rates = append(rates, entities.CurrencyRate{
			Code:  r.CharCode,
			Name:  r.Name,
			Value: value,
			Date:  dateStr,
		})
	}

	return rates, nil
}