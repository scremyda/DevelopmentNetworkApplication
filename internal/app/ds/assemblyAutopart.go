package ds

type AssemblyAutopart struct {
	ID         uint     `json:"id" gorm:"primary_key"`
	AutopartID uint     `json:"autopart_id"`
	AssemblyID uint     `json:"assembly_id"`
	Assemblies Assembly `json:"assembly" gorm:"foreignKey:AssemblyID"`
	Autopart   Autopart `json:"autopart" gorm:"foreignKey:AutopartID"`
	Count      int      `json:"count" gorm:"int"`
}

type AssemblyAutopartUpdate struct {
	ID    uint `json:"id" gorm:"primary_key"`
	Count int  `json:"count" gorm:"int"`
}
