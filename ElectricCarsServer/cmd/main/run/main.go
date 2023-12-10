package main

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/Minio"
	"ElectricCarsServer/ElectricCarsServer/internal/app/config"
	"ElectricCarsServer/ElectricCarsServer/internal/app/dsn"
	"ElectricCarsServer/ElectricCarsServer/internal/app/handlers"
	"ElectricCarsServer/ElectricCarsServer/internal/app/pkg"
	"ElectricCarsServer/ElectricCarsServer/internal/app/pkg/auth"
	"ElectricCarsServer/ElectricCarsServer/internal/app/redis"
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	minioClient := Minio.NewMinioClient(logger)
	router := gin.Default()

	router.Use(corsMiddleware())

	conf, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalf("Error with configuration reading: %s", err)
	}

	ctx := context.Background()
	redisClient, errRedis := redis.New(ctx, conf.Redis)
	if errRedis != nil {
		logger.Fatalf("Errof with redis connect: %s", err)
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

	tokenManager, err := auth.NewManager(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		logger.Fatalln(err)
	}
	hand := handlers.NewHandler(logger, rep, minioClient, conf, redisClient, tokenManager)
	application := pkg.NewApp(conf, router, logger, hand)
	application.RunApp()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
