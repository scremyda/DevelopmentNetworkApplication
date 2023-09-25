package main

import (
	"backened/internal/app/ds"
	"backened/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func main() {
	_ = godotenv.Load()
	env, err2 := dsn.FromEnv()
	if err2 != nil {
		panic("Error from reading env")
	}
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(
		&ds.Users{},
		&ds.AssemblyRequest{},
		&ds.Autopart{},
		&ds.AutopartRequest{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}

	for i := 1; i <= 5; i++ {
		user := ds.Users{Login: "user" + strconv.Itoa(i), Password: "password" + strconv.Itoa(i)}
		db.Create(&user)

		autopart := ds.Autopart{Name: "Autopart" + strconv.Itoa(i), Description: "Description" + strconv.Itoa(i), Brand: "Brand" + strconv.Itoa(i), Models: "Model" + strconv.Itoa(i), Year: 2023, Image: "image" + strconv.Itoa(i) + ".jpg", IsDelete: false}
		db.Create(&autopart)

		assemblyRequest := ds.AssemblyRequest{DateStart: time.Now(), DateEnd: time.Now(), Status: "Pending", Factory: "Factory" + strconv.Itoa(i), UserID: user.ID}
		db.Create(&assemblyRequest)

		autopartRequest := ds.AutopartRequest{ASID: assemblyRequest.ID, AutopartID: autopart.ID}
		db.Create(&autopartRequest)
	}
}
