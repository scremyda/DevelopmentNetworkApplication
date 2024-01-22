package ds

import (
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
	DiscussionWithSupplier string    `gorm:"type:text" json:"discussion"`
}
