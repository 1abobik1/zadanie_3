package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/1abobik1/zadanie_3/internal/entities"
	"github.com/1abobik1/zadanie_3/internal/repository"
	"github.com/sirupsen/logrus"
)

type CurrencyUseCase interface {
	AnalyzeRates() (*entities.AnalysisResult, error)
}

type CurrencyAnalysisUseCase struct {
	repo repository.CurrencyRepository
}

func NewCurrencyAnalysisUseCase(repo repository.CurrencyRepository) *CurrencyAnalysisUseCase {
	return &CurrencyAnalysisUseCase{repo: repo}
}

func (uc *CurrencyAnalysisUseCase) AnalyzeRates() (*entities.AnalysisResult, error) {
	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, 0, -90)

	var (
		mu       sync.Mutex
		allRates []entities.CurrencyRate
		wg       sync.WaitGroup
	)

	concurrency := 10
	sem := make(chan struct{}, concurrency)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {

		wg.Add(1)
		sem <- struct{}{}
		go func(date time.Time) {
			defer wg.Done()
			defer func() { <-sem }()

			dateReq := date.Format("02/01/2006")
			rates, err := uc.repo.GetRates(dateReq)
			if err != nil {
				logrus.Errorf("failed to get rates for %s: %v", dateReq, err)
				return
			}

			mu.Lock()
			allRates = append(allRates, rates...)
			mu.Unlock()
		}(d)
	}

	wg.Wait()

	if len(allRates) == 0 {
		return nil, fmt.Errorf("no currency rates found for the period")
	}

	maxRate := allRates[0]
	minRate := allRates[0]
	total := 0.0

	for _, rate := range allRates {
		if rate.Value > maxRate.Value {
			maxRate = rate
		}
		if rate.Value < minRate.Value {
			minRate = rate
		}
		total += rate.Value
	}

	average := total / float64(len(allRates))

	return &entities.AnalysisResult{
		MaxRate:     maxRate,
		MinRate:     minRate,
		AverageRate: average,
	}, nil
}
