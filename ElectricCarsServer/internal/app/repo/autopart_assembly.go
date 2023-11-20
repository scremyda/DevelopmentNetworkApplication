package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"fmt"
)

func (r *Repository) DeleteFromAssembly(deleteFromAssembly ds.Autopart_Assembly) error {
	var deletedAutopartAssembly ds.Autopart_Assembly
	result := r.db.Where("autopart_id = ? and assembly_id = ?", deleteFromAssembly.AutopartID,
		deleteFromAssembly.AssemblyID).Find(&deletedAutopartAssembly)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	if err := r.db.Delete(&deletedAutopartAssembly).Error; err != nil {
		return err
	}

	return result.Error
}

func (r *Repository) UpdateCountAutopartAssembly(updatedAssembly ds.Autopart_Assembly) error {
	oldAssembly := ds.Autopart_Assembly{}
	result := r.db.Where("autopart_id = ? and assembly_id = ?", updatedAssembly.AutopartID,
		updatedAssembly.AssemblyID).Find(&oldAssembly)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	oldAssembly.Count = updatedAssembly.Count

	result = r.db.Save(oldAssembly)
	return result.Error
}
