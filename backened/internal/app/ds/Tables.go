package ds

import (
	"gorm.io/gorm"
	"time"
)

//// Пользователи (Users)
//type User struct {
//	gorm.Model
//	ID       uint      `gorm:"primaryKey" json:"user_id"`
//	Username string    `gorm:"not null" json:"username"`
//	Email    string    `gorm:"not null" json:"email"`
//	Requests []Request `json:"requests"`
//}
//
//// Услуги (Services)
//type Service struct {
//	gorm.Model
//	ID     uint   `gorm:"primaryKey" json:"service_id"`
//	Name   string `gorm:"not null" json:"name"`
//	Status string `gorm:"not null;check:status IN ('Действует', 'Удалена')" json:"status"`
//}
//
//// Заявки (Requests)
//type Request struct {
//	gorm.Model
//	ID             uint       `gorm:"primaryKey" json:"request_id"`
//	Status         string     `gorm:"not null;check:status IN ('Введена', 'В работе', 'Завершена', 'Отменена', 'Удалена')" json:"status"`
//	CreationDate   *time.Time `gorm:"not null" json:"creation_date"`
//	FormationDate  *time.Time `gorm:"not null" json:"formation_date"`
//	CompletionDate *time.Time `json:"completion_date"`
//	ModeratorID    uint       `json:"moderator_id"`
//	UserID         uint       `gorm:"not null" json:"user_id"`
//	Services       []Service  `gorm:"many2many:request_services;" json:"services"`
//}
//
//// Многие-ко-многим "Заявки-Услуги" (RequestServices)
//type RequestService struct {
//	gorm.Model
//	RequestID uint `gorm:"primaryKey" json:"-"`
//	ServiceID uint `gorm:"primaryKey" json:"-"`
//}
//
//// Статусы Заявок (RequestStatuses)
//type RequestStatus struct {
//	gorm.Model
//	ID   uint   `gorm:"primaryKey" json:"status_id"`
//	Name string `gorm:"not null" json:"name"`
//}

type Users struct {
	gorm.Model
	Login       string `gorm:"type:varchar(255);unique" json:"login"`
	Password    string `gorm:"type:varchar(255)" json:"-"`
	IsModerator bool
}
type Autopart struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Brand       string `gorm:"not null" json:"brand"`
	Models      string `gorm:"not null" json:"model"`
	Year        int    `gorm:"not null" json:"year"`
	Image       string `gorm:"type:varchar(255)" json:"image"`
	IsDelete    bool   `json:"is_delete"`
	UserID      uint   `json:"-"`
	User        Users  `gorm:"foreignKey:UserID" json:"-"`
	Status      string `gorm:"type:varchar(255)" json:"status"`
	Price       uint   `gorm:"not null" json:"price"`
}
type Autopart_Assembly struct {
	gorm.Model
	AssemblyID      uint     `json:"-"`
	AutopartID      uint     `json:"-"`
	AssemblyRequest Assembly `gorm:"foreignKey:AssemblyID" json:"-"`
	Autopart        Autopart `gorm:"foreignKey:AutopartID" json:"-"`
	Cash            uint     `gorm:"not null" json:"cash"`
}

type Assembly struct {
	gorm.Model
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Status      string    `gorm:"type:varchar(255)" json:"status"`
	Name        string    `gorm:"type:varchar(255)" json:"factory"`
	ImageURL    string    `gorm:"type:varchar(255)" json:"image"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}
