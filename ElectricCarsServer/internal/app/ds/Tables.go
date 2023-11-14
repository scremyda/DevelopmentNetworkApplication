package ds

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model
	Login       string `gorm:"type:varchar(255);unique" json:"login"`
	Password    string `gorm:"type:varchar(255)" json:"-"`
	IsModerator bool
}
type Autopart struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(255)" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Brand       string  `gorm:"not null" json:"brand"`
	Models      string  `gorm:"not null" json:"model"`
	Year        int     `gorm:"not null" json:"year"`
	Image       string  `gorm:"type:text" json:"image"`
	UserID      uint    `json:"user_id"`
	User        Users   `gorm:"foreignKey:UserID" json:"-"`
	Status      bool    `gorm:"type:bool" json:"is_deleted"`
	Price       float64 `gorm:"not null" json:"price"`
}
type Autopart_Assembly struct {
	gorm.Model
	AssemblyID      uint     `json:"-"`
	AutopartID      uint     `json:"-"`
	AssemblyRequest Assembly `gorm:"foreignKey:AssemblyID" json:"-"`
	Autopart        Autopart `gorm:"foreignKey:AutopartID" json:"-"`
	Count           int      `json:"count"`
}

type Assembly struct {
	gorm.Model
	DateStart             time.Time `json:"date_start"`
	DateEnd               time.Time `json:"date_end"`
	DateStartOfProcessing time.Time `json:"date_processing"`
	Status                string    `gorm:"type:text" json:"status"`
	Name                  string    `gorm:"type:text" json:"factory"`
	Creator               uint      `json:"creator_id"`
	Description           string    `gorm:"type:text" json:"description"`
}

type AssemblyDetails struct {
	Assembly  *Assembly
	Autoparts *[]Autopart
}

type AutopartDetails struct {
	Autopart_id   int    `json:"autopart_id"`
	Autopart_name string `json:"name"`
}

type AddToAssemblyID struct {
	AutopartDetails AutopartDetails `json:"autopart"`
	Assembly        Assembly        `json:"assembly"`
}
