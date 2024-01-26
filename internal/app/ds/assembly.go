package ds

import (
	"time"
)

// ЗАЯВКА ТЕНДЕРА
type Assembly struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	Name           string    `json:"assembly_name" gorm:"type:varchar(255)"`
	Status         string    `json:"status" gorm:"type:varchar(15)"`
	CreationDate   time.Time `json:"creation_date" gorm:"not null; default:current_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	//CreatorLogin    string          `json:"creator_login"`
	//ModeratorLogin  string          `json:"moderator_login"`
	StatusCheck       string             `json:"status_check"`
	ModeratorID       *uint              `json:"moderator_id"`
	Moderator         User               `json:"moderator" gorm:"foreignkey:ModeratorID"`
	UserID            uint               `json:"user_id"`
	User              User               `json:"user" gorm:"foreignkey:UserID"`
	AssemblyAutoparts []AssemblyAutopart `json:"autopart_assemblies" gorm:"foreignkey:AssemblyID"`
}

type AssemblyResponse struct {
	ID                uint               `json:"id" gorm:"primary_key"`
	Name              string             `json:"assembly_name" gorm:"type:varchar(255)"`
	Status            string             `json:"status" gorm:"type:varchar(15)"`
	CreationDate      time.Time          `json:"creation_date" gorm:"type:datetime; not null; default:current_date"`
	FormationDate     time.Time          `json:"formation_date" gorm:"type:datetime"`
	CompletionDate    time.Time          `json:"completion_date" gorm:"type:datetime"`
	StatusCheck       string             `json:"status_check"`
	AssemblyAutoparts []AssemblyAutopart `json:"autopart_assemblies"`
	UserName          string             `json:"user_name"`
	ModeratorName     string             `json:"moderator_name"`
	UserLogin         string             `json:"user_login"`
	ModeratorLogin    string             `json:"moderator_login"`
}

type NewStatus struct {
	Status     string `json:"status"`
	AssemblyID uint   `json:"assembly_id"`
}

type AssemblyDetails struct {
	Assembly *Assembly   `json:"assembly"`
	Autopart *[]Autopart `json:"autoparts"`
}

type RequestAsyncService struct {
	RequestId uint   `gorm:"primaryKey" json:"requestId"`
	Token     string `json:"Server_Token"`
	Status    string `json:"status"`
}

type UpdateAssembly struct {
	ID   uint   `json:"id"`
	Name string `json:"assembly_name"`
}
