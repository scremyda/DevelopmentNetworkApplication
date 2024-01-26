package repository

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) AssemblyByUserID(userID string) (*[]ds.AssemblyResponse, error) {
	var assemblies []ds.Assembly
	var assemblyResponses = []ds.AssemblyResponse{}
	result := r.db.Preload("User").
		Preload("Moderator").
		Where("user_id = ? AND status != 'удален' AND status != 'черновик'", userID).
		Find(&assemblies)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, tender := range assemblies {
		tenderResponse := ds.AssemblyResponse{
			ID:        tender.ID,
			Name:      tender.Name,
			UserName:  tender.User.Name,
			UserLogin: tender.User.Login,
			//UserRole:       tender.User.Role,
			ModeratorName: tender.Moderator.Name,
			Status:        tender.Status,
			StatusCheck:   tender.StatusCheck,
			//ModeratorRole:  tender.Moderator.Role,
			ModeratorLogin:    tender.Moderator.Login,
			CreationDate:      tender.CreationDate,
			FormationDate:     tender.FormationDate,
			CompletionDate:    tender.CompletionDate,
			AssemblyAutoparts: tender.AssemblyAutoparts,
		}
		assemblyResponses = append(assemblyResponses, tenderResponse)
	}

	return &assemblyResponses, result.Error
}

func (r *Repository) AssemblyByID(id uint) (*ds.AssemblyResponse, error) {
	assembly := ds.Assembly{}

	result := r.db.Preload("User").
		//Preload("TenderCompanies.Tenders").
		Preload("AssemblyAutoparts.Autopart").
		First(&assembly, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	assemblyResponse := ds.AssemblyResponse{
		ID:                assembly.ID,
		Name:              assembly.Name,
		UserName:          assembly.User.Name,
		UserLogin:         assembly.User.Login,
		ModeratorName:     assembly.Moderator.Name,
		Status:            assembly.Status,
		StatusCheck:       assembly.StatusCheck,
		ModeratorLogin:    assembly.Moderator.Login,
		CreationDate:      assembly.CreationDate,
		FormationDate:     assembly.FormationDate,
		CompletionDate:    assembly.CompletionDate,
		AssemblyAutoparts: assembly.AssemblyAutoparts,
	}
	return &assemblyResponse, result.Error
}

func (r *Repository) AssemblyModel(id uint) (*ds.Assembly, error) {
	assembly := ds.Assembly{}

	result := r.db.Preload("User").
		//Preload("TenderCompanies.Tenders").
		Preload("AssemblyAutoparts.Autopart").
		First(&assembly, id)
	return &assembly, result.Error
}

func (r *Repository) AssemblyDraftId(userId uint) (uint, error) {
	var assembly ds.Assembly
	result := r.db.
		Where("status = ? AND user_id = ?", "черновик", userId).
		First(&assembly)
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return assembly.ID, result.Error
}

func (r *Repository) GetAssemblyDraftID(creatorID uint) (uint, error) {
	var draftReq ds.Assembly

	res := r.db.Where("user_id = ?", creatorID).Where("status = ?", utils.Draft).Take(&draftReq)
	if errors.Is(gorm.ErrRecordNotFound, res.Error) {
		return 0, nil
	}

	if res.Error != nil {
		return 0, res.Error
	}

	return draftReq.ID, nil
}

func (r *Repository) CreateAssemblyDraft(creatorID uint) (uint, error) {
	request := ds.Assembly{
		UserID:       creatorID,
		Status:       "черновик",
		CreationDate: time.Now(),
		ModeratorID:  nil,
	}

	if err := r.db.Create(&request).Error; err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (r *Repository) GetAssemblyWithDataByID(requestID uint, userId uint, isAdmin bool) (ds.Assembly, []ds.Autopart, error) {
	var AssemblyRequest ds.Assembly
	var autoparts []ds.Autopart

	//ищем такую заявку
	result := r.db.First(&AssemblyRequest, "id =?", requestID)
	if result.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return ds.Assembly{}, nil, result.Error
	}
	if !isAdmin && AssemblyRequest.UserID == uint(userId) || isAdmin {
		res := r.db.
			Table("assembly_autoparts").
			Select("autoparts.*").
			Where("status != ?", "удалён").
			Joins("JOIN autoparts ON assembly_autoparts.\"AutopartID\" = autoparts.id").
			Where("assembly_autoparts.\"AssemblyID\" = ?", requestID).
			Find(&autoparts)
		if res.Error != nil {
			r.logger.Error("error while getting for assembly request")
			return ds.Assembly{}, nil, res.Error
		}
	} else {
		return ds.Assembly{}, nil, errors.New("ошибка доступа к данной заявке")
	}

	return AssemblyRequest, autoparts, nil
}

func (r *Repository) AssembliesList(statusID string, startDate time.Time, endDate time.Time) (*[]ds.AssemblyResponse, error) {
	var assemblies []ds.Assembly
	assemblyResponses := []ds.AssemblyResponse{}
	if statusID == "" {
		result := r.db.
			Preload("User").
			Preload("Moderator").
			Where("status != 'удален' AND status != 'черновик' AND creation_date BETWEEN ? AND ?", startDate, endDate).
			Find(&assemblies)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, nil
			} else {
				return nil, result.Error
			}
		}

		for _, assembly := range assemblies {
			assemblyResponse := ds.AssemblyResponse{
				ID:        assembly.ID,
				Name:      assembly.Name,
				UserName:  assembly.User.Name,
				UserLogin: assembly.User.Login,
				//UserRole:       tender.User.Role,
				ModeratorName: assembly.Moderator.Name,
				Status:        assembly.Status,
				StatusCheck:   assembly.StatusCheck,
				//ModeratorRole:  tender.Moderator.Role,
				ModeratorLogin:    assembly.Moderator.Login,
				CreationDate:      assembly.CreationDate,
				FormationDate:     assembly.FormationDate,
				CompletionDate:    assembly.CompletionDate,
				AssemblyAutoparts: assembly.AssemblyAutoparts,
			}
			assemblyResponses = append(assemblyResponses, assemblyResponse)
		}

		return &assemblyResponses, result.Error
	}

	result := r.db.
		Preload("User").
		Where("status = ? AND status != 'черновик' AND creation_date BETWEEN ? AND ?", statusID, startDate, endDate).
		Find(&assemblies)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	for _, assembly := range assemblies {
		assemblyResponse := ds.AssemblyResponse{
			ID:        assembly.ID,
			Name:      assembly.Name,
			UserName:  assembly.User.Name,
			UserLogin: assembly.User.Login,
			Status:    assembly.Status,
			//UserRole:       tender.User.Role,
			ModeratorName: assembly.Moderator.Name,
			//ModeratorRole:  tender.Moderator.Role,
			ModeratorLogin:    assembly.Moderator.Login,
			StatusCheck:       assembly.StatusCheck,
			CreationDate:      assembly.CreationDate,
			FormationDate:     assembly.FormationDate,
			CompletionDate:    assembly.CompletionDate,
			AssemblyAutoparts: assembly.AssemblyAutoparts,
		}
		assemblyResponses = append(assemblyResponses, assemblyResponse)
	}

	return &assemblyResponses, result.Error
}

func (r *Repository) UpdateAssembly(updatedAssembly *ds.Assembly) error {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, updatedAssembly.ID); result.Error != nil {
		return result.Error
	}
	if updatedAssembly.Name != "" {
		oldAssembly.Name = updatedAssembly.Name
	}
	if updatedAssembly.CreationDate.String() != utils.EmptyDate {
		oldAssembly.CreationDate = updatedAssembly.CreationDate
	}
	if updatedAssembly.CompletionDate.String() != utils.EmptyDate {
		oldAssembly.CompletionDate = updatedAssembly.CompletionDate
	}
	if updatedAssembly.FormationDate.String() != utils.EmptyDate {
		oldAssembly.FormationDate = updatedAssembly.FormationDate
	}
	if updatedAssembly.Status != "" {
		oldAssembly.Status = updatedAssembly.Status
	}

	*updatedAssembly = oldAssembly
	result := r.db.Save(updatedAssembly)
	return result.Error
}

func (r *Repository) FormAssemblyRequestByIDAsynce(id uint, creatorID uint) (error, uint) {
	var req ds.Assembly
	res := r.db.
		Where("id = ?", id).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, 0
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), 0
	}

	req.StatusCheck = "В обработке"
	req.Status = "сформирован"
	req.FormationDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err, 0
	}

	return nil, req.ID
}

func (r *Repository) FormAssemblyRequestByID(creatorID uint) (error, uint) {
	var req ds.Assembly
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, 0
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), 0
	}
	req.StatusCheck = "В обработке"
	req.Status = "сформирован"
	req.FormationDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err, 0
	}

	return nil, req.ID
}

func (r *Repository) GetAssemblyByUser(creatorID uint) (error, uint) {
	var req ds.Assembly
	res := r.db.
		//Where("id = ?", requestID).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, 0
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), 0
	}

	return nil, req.ID
}

func (r *Repository) GetAssemblyByID(creatorID uint, id uint) (error, ds.Assembly) {
	var req ds.Assembly
	res := r.db.
		Where("id = ?", id).
		Where("user_id = ?", creatorID).
		Where("status = ?", utils.Draft).
		Take(&req)

	if res.Error != nil {
		return res.Error, req
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки"), req
	}

	return nil, req
}

func (r *Repository) FinishRejectHelper(status string, requestID, moderatorID uint) error {
	var req ds.Assembly
	res := r.db.
		Where("id = ?", requestID).
		Where("status = ?", "сформирован").
		Take(&req)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("нет такой заявки")
	}

	req.ModeratorID = &moderatorID
	//req.ModeratorLogin = userInfo.Login
	req.Status = status

	req.CompletionDate = time.Now()

	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAssemblyByID(requestID uint) error { // ?
	var req ds.Assembly
	if result := r.db.First(&req, requestID); result.Error != nil {
		return result.Error
	}

	req.Status = "удален"
	req.CompletionDate = time.Now()
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	result := r.db.Save(&req)

	return result.Error
}

func (r *Repository) DeleteAutopartFromRequest(id int) error {
	var dh ds.AssemblyAutopart
	if result := r.db.First(&dh, id); result.Error != nil {
		return result.Error
	}
	return r.db.Delete(&dh).Error
}

func (r *Repository) UpdateAssemblyAutopart(id uint, count int) error {
	var updateAutopart ds.AssemblyAutopart
	r.db.Where("id = ?", id).First(&updateAutopart)

	if updateAutopart.AssemblyID == 0 {
		return errors.New("нет такой заявки")
	}
	updateAutopart.Count = count

	if err := r.db.Save(&updateAutopart).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveRequest(monitoringRequest ds.RequestAsyncService) error {
	var request ds.Assembly
	err := r.db.First(&request, "id = ?", monitoringRequest.RequestId)
	if err.Error != nil {
		r.logger.Error("error while getting monitoring request")
		return err.Error
	}

	request.StatusCheck = monitoringRequest.Status
	res := r.db.Save(&request)
	return res.Error
}
