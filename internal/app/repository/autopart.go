package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"strings"
)

func (r *Repository) GetOpenAutoparts() (*[]ds.Autopart, error) {
	var tenders []ds.Autopart
	if err := r.db.Where("status = ?", "действует").Find(&tenders).Error; err != nil {
		return nil, err
	}
	return &tenders, nil
}

func (r *Repository) GetAutopartById(id uint) (*ds.Autopart, error) {
	var company ds.Autopart
	if err := r.db.Where("status = ?", "действует").First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *Repository) AutopartsList(name string) (*[]ds.Autopart, error) {
	var autoparts []ds.Autopart
	if err := r.db.Where("autopart_name LIKE ? AND status != ?", "%"+name+"%", "удален").Find(&autoparts).Error; err != nil {
		return nil, err
	}
	return &autoparts, nil
}

func (r *Repository) AddAutopart(autopart *ds.Autopart) (uint, error) {
	autopart.Status = "действует"
	result := r.db.Create(&autopart)
	return autopart.ID, result.Error
}

func (r *Repository) DeleteAutopart(id uint) error {
	autopart := ds.Autopart{}

	if err := r.db.First(&autopart, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&autopart).Update("status", "удален").Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAutopart(updatedAutopart *ds.Autopart) error {
	var oldAutopart ds.Autopart

	if result := r.db.First(&oldAutopart, updatedAutopart.ID); result.Error != nil {
		return result.Error
	}

	if updatedAutopart.AutopartName != "" {
		oldAutopart.AutopartName = updatedAutopart.AutopartName
	}

	if updatedAutopart.ImageURL != "" {
		oldAutopart.ImageURL = updatedAutopart.ImageURL
	}

	if updatedAutopart.Description != "" {
		oldAutopart.Description = updatedAutopart.Description
	}

	if updatedAutopart.Status != "" {
		oldAutopart.Status = updatedAutopart.Status
	}

	if updatedAutopart.Price != 0 {
		oldAutopart.Price = updatedAutopart.Price
	}

	if updatedAutopart.Year != 0 {
		oldAutopart.Year = updatedAutopart.Year
	}

	*updatedAutopart = oldAutopart
	result := r.db.Save(updatedAutopart)
	return result.Error
}

func (r *Repository) DeleteAutopartImage(companyId uint) string {
	company := ds.Autopart{}

	r.db.First(&company, "id = ?", companyId)
	return company.ImageURL
}

func (r *Repository) AddAutopartToDraft(dataID uint, creatorID uint, count int) (uint, error) {
	// получаем услугу
	data, err := r.GetAutopartById(dataID)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("нет такой услуги")
	}
	if data.Status == "удален" {
		return 0, errors.New("услуга удалена")
	}

	// получаем черновик
	var draftReq ds.Assembly
	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", utils.Draft).Take(&draftReq)

	// создаем черновик, если его нет
	if res.RowsAffected == 0 {
		newDraftRequestID, err := r.CreateAssemblyDraft(creatorID)
		if err != nil {
			return 0, err
		}

		draftReq.ID = newDraftRequestID
	}

	// добавляем запись в мм
	requestToData := ds.AssemblyAutopart{
		AutopartID: dataID,
		AssemblyID: draftReq.ID,
		Count:      count,
	}

	err = r.db.Create(&requestToData).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return 0, errors.New("услуга уже существует в заявке")
		}

		return 0, err
	}

	return draftReq.ID, nil
}
