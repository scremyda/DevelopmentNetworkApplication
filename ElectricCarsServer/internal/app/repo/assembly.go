package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/utils"
	"fmt"
	"gorm.io/gorm/clause"
)

func (r *Repository) AssembliesList() (*[]ds.Assembly, error) {
	var assemblies []ds.Assembly
	result := r.db.Where("status != ?", utils.DeletedString).Find(&assemblies)
	return &assemblies, result.Error
}

func (r *Repository) AssemblyByID(id uint) (*[]ds.Autopart, *ds.Assembly, error) {
	var autoparts []ds.Autopart
	var assembly ds.Assembly

	result := r.db.Where("status != ? and id = ?", utils.DeletedString, id).Find(&assembly)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	resultAutoparts := r.db.Preload(clause.Associations).
		Joins("JOIN autopart_assemblies ON autopart_assemblies.autopart_id = autoparts.id").
		Joins("JOIN assemblies ON autopart_assemblies.assembly_id = assemblies.id").
		Where("assemblies.id = ?", id).
		Find(&autoparts)

	if resultAutoparts.Error != nil {
		return nil, nil, resultAutoparts.Error
	}

	return &autoparts, &assembly, nil
}

//func (r *Repository) AddAssembly(assembly *ds.Assembly) error {
//	result := r.db.Create(&assembly)
//	return result.Error
//}

func (r *Repository) DeleteAssembly(id uint) error {
	var assembly ds.Assembly
	if result := r.db.First(&assembly, id); result.Error != nil {
		return result.Error
	}
	if assembly.ID == 0 {
		return fmt.Errorf("assembly not found")
	}
	assembly.Status = utils.DeletedString
	result := r.db.Save(&assembly)
	return result.Error
}

func (r *Repository) UpdateAssembly(updatedAssembly *ds.Assembly) error {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, updatedAssembly.ID); result.Error != nil {
		return result.Error
	}
	if updatedAssembly.Name != "" {
		oldAssembly.Name = updatedAssembly.Name
	}
	if updatedAssembly.DateStart.String() != utils.EmptyDate {
		oldAssembly.DateStart = updatedAssembly.DateStart
	}
	if updatedAssembly.DateEnd.String() != utils.EmptyDate {
		oldAssembly.DateEnd = updatedAssembly.DateEnd
	}
	if updatedAssembly.DateStartOfProcessing.String() != utils.EmptyDate {
		oldAssembly.DateStartOfProcessing = updatedAssembly.DateStartOfProcessing
	}
	if updatedAssembly.Status != "" {
		oldAssembly.Status = updatedAssembly.Status
	}
	if updatedAssembly.Description != "" {
		oldAssembly.Description = updatedAssembly.Description
	}
	*updatedAssembly = oldAssembly
	result := r.db.Save(updatedAssembly)
	return result.Error
}
