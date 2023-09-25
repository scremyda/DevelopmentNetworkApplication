package main

import (
	"backened/internal/app/config"
	"backened/internal/app/dsn"
	"backened/internal/app/handler"
	"backened/internal/app/pkg"
	"backened/internal/app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	router := gin.Default()
	conf, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalf("Error with configuration reading: %s", err)
	}
	postgresString, errPost := dsn.FromEnv()
	if errPost != nil {
		logger.Fatalf("Error of reading postgres line: %s", errPost)
	}
	fmt.Println(postgresString)
	rep, errRep := repository.NewRepository(postgresString, logger)
	if errRep != nil {
		logger.Fatalf("Error from repository: %s", err)
	}
	hand := handler.NewHandler(logger, rep)
	application := pkg.NewApp(conf, router, logger, hand)
	application.RunApp()
}
