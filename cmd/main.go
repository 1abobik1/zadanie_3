package main

import (
	"github.com/1abobik1/zadanie_3/internal/delivery/http/handler"
	"github.com/1abobik1/zadanie_3/internal/repository"
	"github.com/1abobik1/zadanie_3/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	
	r := gin.Default()

	repo := repository.NewCBRApiClient()
	uc := usecase.NewCurrencyAnalysisUseCase(repo)
	handler := handler.NewCurrencyHandler(uc)

	r.GET("/analysis", handler.GetAnalysis)

	r.Run(":8080")
}
