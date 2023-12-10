package main

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		&ds.Assembly{},
		&ds.Autopart{},
		&ds.Autopart_Assembly{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}
	users := []ds.Users{
		{Login: "user1", Name: "user1", Password: "password1", IsModerator: true},
		{Login: "user2", Name: "user2", Password: "password2", IsModerator: false},
		{Login: "user3", Name: "user3", Password: "password3", IsModerator: false},
	}
	db.Create(&users)
}
