package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/1abobik1/zadanie_3/internal/usecase"
)

type CurrencyHandler struct {
	useCase usecase.CurrencyUseCase
}

func NewCurrencyHandler(uc usecase.CurrencyUseCase) *CurrencyHandler {
	return &CurrencyHandler{uc}
}

func (h *CurrencyHandler) GetAnalysis(c *gin.Context) {
	startTime := time.Now()
	
	logrus.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})

	logFields := logrus.Fields{
		"duration": time.Since(startTime).String(),
		"status":   c.Writer.Status(),
	}

	result, err := h.useCase.AnalyzeRates()
	if err != nil {
        logrus.WithFields(logFields).WithError(err).Error("Request processing failed")
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


	c.JSON(http.StatusOK, result)
}
