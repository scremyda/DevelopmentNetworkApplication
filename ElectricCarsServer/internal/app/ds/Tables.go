package ds

import (
	"time"
)

type Users struct {
	ID          uint   `gorm:"primary_key" json:"user_id"`
	Login       string `gorm:"type:varchar(255);unique" json:"login"`
	Password    string `gorm:"type:varchar(255)" json:"password"`
	IsModerator bool   `json:"is_moderator"`
}
type Autopart struct {
	ID          uint    `gorm:"primary_key" json:"autopart_id"`
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
	ID              uint     `gorm:"primary_key" json:"autopart_factory_id"`
	AssemblyID      uint     `json:"factory_id"`
	AutopartID      uint     `json:"autopart_id"`
	AssemblyRequest Assembly `gorm:"foreignKey:AssemblyID" json:"-"`
	Autopart        Autopart `gorm:"foreignKey:AutopartID" json:"-"`
	Count           int      `json:"count"`
}

type Assembly struct {
	ID                    uint      `gorm:"primary_key" json:"factory_id"`
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
	User_id         uint            `json:"user_id"`
}

type AssemblyForm struct {
	Factory_id uint `json:"factory_id"`
	User_id    uint `json:"user_id"`
}
