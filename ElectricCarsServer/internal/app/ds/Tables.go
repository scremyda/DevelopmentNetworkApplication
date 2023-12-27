package ds

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Users struct {
	ID          uint   `gorm:"primary_key" json:"user_id"`
	Login       string `gorm:"type:varchar(255);unique" json:"login"`
	Name        string `json:"name"`
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
	ID                     uint      `gorm:"primary_key" json:"factory_id"`
	DateStart              time.Time `json:"date_start"`
	DateEnd                time.Time `json:"date_end"`
	DateStartOfProcessing  time.Time `json:"date_processing"`
	Status                 string    `gorm:"type:text" json:"status"`
	Name                   string    `gorm:"type:text" json:"factory"`
	Creator                uint      `json:"creator_id"`
	Description            string    `gorm:"type:text" json:"description"`
	CreatorLogin           string    `json:"creator_login"`
	AdminLogin             string    `json:"admin_login"`
	DiscussionWithSupplier string    `gorm:"type:text" json:"discussion"`
}

type AssemblyAdmin struct {
	Assembly   Assembly
	AdminLogin string `json:"admin_login"`
}

type AssemblyDetails struct {
	Assembly  *Assembly
	Autoparts *[]Autopart
}

type AutopartList struct {
	DraftID   int         `json:"draft_id"`
	Autoparts *[]Autopart `json:"autoparts_list"`
}

type AddToAssemblyID struct {
	AutopartDetails AutopartDetails `json:"autopart"`
	User_id         uint            `json:"user_id"`
}

type AssemblyForm struct {
	Factory_id uint   `json:"factory_id"`
	User_id    uint   `json:"user_id"`
	Status     string `json:"status"`
}

type AutopartDetails struct {
	Autopart_id   int    `json:"autopart_id"`
	Autopart_name string `json:"name"`
}

type JwtClaims struct {
	jwt.StandardClaims
	UserId  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
}

type Role int

const (
	Client Role = iota // 0
	Admin              // 1
)

var (
	ErrClientAlreadyExists = errors.New("клиент с таким логином уже существует")
	ErrUserNotFound        = errors.New("клиента с таким логином не существует")
)

type UserLogin struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=4,max=64"`
}

type UserSignUp struct {
	Login    string `json:"login" binding:"required,max=64"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required,min=4,max=64"`
}

type RequestAsyncService struct {
	AssemblyID             int    `json:"assemblyId"`
	DiscussionWithSupplier string `json:"discussion"`
	Token                  string `json:"Token"`
}
