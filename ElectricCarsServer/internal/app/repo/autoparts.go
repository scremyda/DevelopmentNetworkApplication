package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
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

func (r *Repository) AddToAssembly(autopartDetails *ds.AddToAssemblyID) error {
	var autopart ds.Autopart
	if err := r.db.Where("id = ? AND name = ?", autopartDetails.AutopartDetails.Autopart_id,
		autopartDetails.AutopartDetails.Autopart_name).
		First(&autopart).Error; err != nil {
		return err
	}
	request := ds.Assembly{
		DateStart: time.Now(),
		Status:    utils.DraftString,
		Creator:   autopartDetails.User_id,
	}

	// Проверка наличия записи с Status = utils.DraftString
	var existingAssembly ds.Assembly
	result := r.db.First(&existingAssembly, "creator = ? AND status = ?", autopartDetails.User_id, utils.DraftString)

	if result.Error != nil {
		// Если записи с Status = utils.DraftString у пользователя нет, создаем новую запись
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := r.db.Create(&request).Error; err != nil {
				return err
			}
		} else {
			// Обработка других ошибок базы данных
			return result.Error
		}
	}

	autopartID := uint(autopartDetails.AutopartDetails.Autopart_id)
	assemblyID := request.ID
	if assemblyID == 0 {
		assemblyID = existingAssembly.ID
	}

	// Поиск записи по autopartID и assemblyID
	var autopartAssembly ds.Autopart_Assembly
	result = r.db.First(&autopartAssembly, "autopart_id = ? AND assembly_id = ?", autopartID, assemblyID)

	if result.Error != nil {
		// Если записи нет, создаем новую запись
		autopartAssembly = ds.Autopart_Assembly{
			AutopartID: autopartID,
			AssemblyID: assemblyID,
			Count:      1,
		}

		if err := r.db.Create(&autopartAssembly).Error; err != nil {
			return err
		}
	} else {
		// Если запись существует, увеличиваем Count на 1 и обновляем запись
		autopartAssembly.Count++
		if err := r.db.Save(&autopartAssembly).Error; err != nil {
			return err
		}
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
