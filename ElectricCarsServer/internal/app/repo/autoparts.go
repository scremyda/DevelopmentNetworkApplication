package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"fmt"
	"gorm.io/gorm"
)

func (r *Repository) AutopartsList(brand string) (*[]ds.Autopart, error) {
	var autoparts []ds.Autopart
	var result *gorm.DB
	if brand == "" {
		result = r.db.Where("status = ?", false).Find(&autoparts)
	} else {
		result = r.db.Where("status = ? AND brand = ?", false, brand).Find(&autoparts)
	}
	return &autoparts, result.Error
}

func (r *Repository) AutopartById(id uint) (*ds.Autopart, error) {
	autopart := ds.Autopart{}
	result := r.db.First(&autopart, id)
	return &autopart, result.Error
}

func (r *Repository) UpdateAutopartImage(id string, newImageURL string) error {
	autopart := ds.Autopart{}
	if result := r.db.First(&autopart, id); result.Error != nil {
		return result.Error
	}
	autopart.Image = newImageURL
	result := r.db.Save(autopart)
	return result.Error
}

func (r *Repository) DeleteAutopart(id uint) error {
	var autopart ds.Autopart
	if result := r.db.First(&autopart, id); result.Error != nil {
		return result.Error
	}
	if autopart.ID == 0 {
		return fmt.Errorf("autopart not found")
	}
	err := r.deleteImageFromMinio(autopart.Image)
	if err != nil {
		return err
	}
	autopart.Status = true
	result := r.db.Save(&autopart)
	return result.Error
}

func (r *Repository) AddAutopart(autopart *ds.Autopart) error {
	result := r.db.Create(&autopart)
	return result.Error
}

func (r *Repository) AddToAssembly(autopartDetails *ds.AutopartDetails, assembly *ds.Assembly) error {
	var autopart ds.Autopart
	if err := r.db.Where("id = ? AND name = ?", autopartDetails.Autopart_id, autopartDetails.Autopart_name).
		First(&autopart).Error; err != nil {
		return err
	}
	assembly.Creator = autopart.UserID
	result := r.db.Where("name = ?", assembly.Name).FirstOrCreate(&assembly)
	if result.Error != nil {
		return result.Error
	}
	autopartAssembly := ds.Autopart_Assembly{
		AutopartID: uint(autopartDetails.Autopart_id),
		AssemblyID: assembly.ID,
		Count:      1,
	}

	if err := r.db.Create(&autopartAssembly).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAutopart(updatedAutopart *ds.Autopart) error {
	var oldAutopart ds.Autopart
	if result := r.db.First(&oldAutopart, updatedAutopart.ID); result.Error != nil {
		return result.Error
	}
	if updatedAutopart.Name != "" {
		oldAutopart.Name = updatedAutopart.Name
	}
	if updatedAutopart.Description != "" {
		oldAutopart.Description = updatedAutopart.Description
	}
	if updatedAutopart.Brand != "" {
		oldAutopart.Brand = updatedAutopart.Brand
	}
	if updatedAutopart.Models != "" {
		oldAutopart.Models = updatedAutopart.Models
	}
	if updatedAutopart.Year != 0 {
		oldAutopart.Year = updatedAutopart.Year
	}
	if updatedAutopart.Image != "" {
		oldAutopart.Image = updatedAutopart.Image
	}
	if updatedAutopart.Price != 0 {
		oldAutopart.Description = updatedAutopart.Description
	}
	oldAutopart.Status = updatedAutopart.Status

	*updatedAutopart = oldAutopart
	result := r.db.Save(updatedAutopart)
	return result.Error
}
