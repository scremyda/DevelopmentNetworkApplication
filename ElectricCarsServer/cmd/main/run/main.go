package main

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/Minio"
	"ElectricCarsServer/ElectricCarsServer/internal/app/config"
	"ElectricCarsServer/ElectricCarsServer/internal/app/dsn"
	"ElectricCarsServer/ElectricCarsServer/internal/app/handlers"
	"ElectricCarsServer/ElectricCarsServer/internal/app/pkg"
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	minioClient := Minio.NewMinioClient(logger)
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
	rep, errRep := repo.NewRepository(postgresString, logger, minioClient)
	if errRep != nil {
		logger.Fatalf("Error from repository: %s", err)
	}
	hand := handlers.NewHandler(logger, rep, minioClient)
	application := pkg.NewApp(conf, router, logger, hand)
	application.RunApp()
}
