package ds

// УСЛУГА
type Autopart struct {
	ID           uint    `json:"autopart_id" gorm:"primary_key"`
	AutopartName string  `json:"name" gorm:"type:varchar(30);not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Price        float64 `json:"price" gorm:"type:float"`
	Year         int     `json:"year" gorm:"type:int"`
	Status       string  `json:"status" gorm:"type:varchar(20);not null"`
	ImageURL     string  `json:"image_url" gorm:"type:varchar(500)"`
}

type AutopartList struct {
	DraftID   uint        `json:"draft_id"`
	Autoparts *[]Autopart `json:"autoparts_list"`
}
type AddToAutopartID struct {
	AutopartID uint `json:"autopart_id"`
	Count      int  `json:"count"`
}
