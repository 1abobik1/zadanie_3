package entities

type CurrencyRate struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Date  string  `json:"date"`
}

type AnalysisResult struct {
	MaxRate     CurrencyRate `json:"max_rate"`
	MinRate     CurrencyRate `json:"min_rate"`
	AverageRate float64      `json:"average_rate"`
}
