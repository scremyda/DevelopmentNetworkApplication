package main

import (
	"backened/internal/app/ds"
	"backened/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		&ds.Assembly{},
		&ds.Autopart{},
		&ds.Autopart_Assembly{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}

	users := []ds.Users{
		{Login: "user1", Password: "password1", IsModerator: true},
		{Login: "user2", Password: "password2", IsModerator: false},
		{Login: "user3", Password: "password3", IsModerator: false},
		//{Login: "user4", Password: "password4", IsModerator: false},
		//{Login: "user5", Password: "password5", IsModerator: false},
	}

	autoparts := []ds.Autopart{
		{Name: "Двигатель Tesla Model 3", Description: "Двигатель Tesla Model 3, 7651488", Brand: "Tesla",
			Models: "Tesla Model 3", Year: 2019, Image: "image1.jpg", IsDelete: false, UserID: 1, Status: "Available", Price: 275000},

		{Name: "Двигатель Tesla Model Y", Description: "Задний мотор, ротор(якорь) мотора, Tesla Model 3, Y, 439210", Brand: "Tesla",
			Models: "Tesla Model Y", Year: 2020, Image: "image2.jpg", IsDelete: false, UserID: 2, Status: "Available", Price: 45457},

		{Name: "Двигатель Tesla Model 3", Description: "Задний мотор, статор и ротор (якорь), Tesla Model 3, Y, 112098000C", Brand: "Tesla",
			Models: "Tesla Model 3", Year: 2020, Image: "image3.jpg", IsDelete: false, UserID: 3, Status: "Available", Price: 45700},
	}
	assemblies := []ds.Assembly{
		{DateStart: time.Now(), DateEnd: time.Now().Add(24 * time.Hour), Status: "Pending", Name: "Завод по сборке в Москве", ImageURL: "factory1.jpg", Description: "Завод по сборке в Москве"},
		{DateStart: time.Now(), DateEnd: time.Now().Add(24 * time.Hour), Status: "Pending", Name: "Завод по сборке в Саратове", ImageURL: "factory2.jpg", Description: "Завод по сборке в Саратове"},
		{DateStart: time.Now(), DateEnd: time.Now().Add(24 * time.Hour), Status: "Pending", Name: "Завод по сборке в Владивостоке", ImageURL: "factory3.jpg", Description: "Завод по сборке в Владивостоке"},
	}

	autopartAssemblies := []ds.Autopart_Assembly{
		{AssemblyID: 1, AutopartID: 1, Cash: 120999},
		{AssemblyID: 2, AutopartID: 2, Cash: 278999},
		{AssemblyID: 3, AutopartID: 3, Cash: 57999},
	}

	db.Create(&users)

	db.Create(&autoparts)

	db.Create(&assemblies)

	db.Create(&autopartAssemblies)

}
