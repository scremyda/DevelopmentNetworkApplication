package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/utils"
	"errors"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var (
	NotSameUser      = errors.New("not the same user")
	AssemblyNotFound = errors.New("assembly not found")
	UserNotFound     = errors.New("user not found")
)

func (r *Repository) SaveAssemblyDiscussion(assemblyAsync ds.RequestAsyncService) error {
	var request ds.Assembly
	err := r.db.First(&request, "id = ? AND status != ? AND status != ?", assemblyAsync.AssemblyID,
		utils.DeletedString, utils.DraftString)
	if err.Error != nil {
		r.logger.Error("error while getting assembly async")
		return err.Error
	}
	request.DiscussionWithSupplier = assemblyAsync.DiscussionWithSupplier
	res := r.db.Save(&request)
	return res.Error
}

func (r *Repository) AssembliesList(status, start, end string, userId int, isAdmin bool) (*[]ds.Assembly, error) {
	var assemblies []ds.Assembly
	ending := "AND creator = " + strconv.Itoa(userId)
	if isAdmin {
		ending = ""
	}
	query := r.db.Where("status != ? AND status != ?"+ending, utils.DeletedString, utils.DraftString)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if start != "" {
		query = query.Where("date_start_of_processing >= ?", start)
	}

	if end != "" {
		query = query.Where("date_start_of_processing <= ?", end)
	}
	query = query.Order("id ASC")
	result := query.Find(&assemblies)
	return &assemblies, result.Error
}

func (r *Repository) AssemblyByID(id uint, userId int, isAdmin bool) (*[]ds.Autopart, *ds.Assembly, error) {
	var autoparts []ds.Autopart
	var assembly ds.Assembly

	result := r.db.Where("status != ? and id = ?", utils.DeletedString, id).Find(&assembly)

	if result.Error != nil {
		return nil, nil, result.Error
	}
	if !isAdmin && assembly.Creator == uint(userId) || isAdmin {
		//ищем услуги в заявке
		resultAutoparts := r.db.Preload(clause.Associations).
			Joins("JOIN autopart_assemblies ON autopart_assemblies.autopart_id = autoparts.id").
			Joins("JOIN assemblies ON autopart_assemblies.assembly_id = assemblies.id").
			Where("assemblies.id = ?", id).
			Find(&autoparts)

		if resultAutoparts.Error != nil {
			return nil, nil, resultAutoparts.Error
		}
	} else {
		return nil, nil, errors.New("ошибка доступа к данной заявке")
	}

	return &autoparts, &assembly, nil
}

//func (r *Repository) AddAssembly(assembly *ds.Assembly) error {
//	result := r.db.Create(&assembly)
//	return result.Error
//}

//func (r *Repository) DeleteAssembly(id uint) error {
//	var assembly ds.Assembly
//	if result := r.db.First(&assembly, id); result.Error != nil {
//		return result.Error
//	}
//	if assembly.ID == 0 {
//		return fmt.Errorf("assembly not found")
//	}
//	assembly.Status = utils.DeletedString
//	result := r.db.Save(&assembly)
//	return result.Error
//}

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
	if updatedAssembly.Description != "" {
		oldAssembly.Description = updatedAssembly.Description
	}
	*updatedAssembly = oldAssembly
	result := r.db.Save(updatedAssembly)
	return result.Error
}

func (r *Repository) UserAdmin(userID uint) (bool, error) {
	var user ds.Users

	result := r.db.First(&user, userID)
	if result.Error != nil {
		return false, result.Error
	}

	if user.ID == 0 {
		return false, UserNotFound
	}

	return user.IsModerator, nil
}

func (r *Repository) SameUser(userID uint, factoryID uint) error {
	var assembly ds.Assembly

	result := r.db.First(&assembly, factoryID)
	if result.Error != nil {
		return result.Error
	}

	if userID != assembly.Creator {
		return NotSameUser
	}

	return nil
}

func (r *Repository) DeleteAssembly(formAssembly ds.AssemblyForm) (ds.Assembly, error) {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, formAssembly.Factory_id); result.Error != nil {
		return ds.Assembly{}, result.Error
	}

	if oldAssembly.ID == 0 {
		return ds.Assembly{}, AssemblyNotFound
	}

	oldAssembly.Status = utils.DeletedString
	oldAssembly.DateEnd = time.Now()

	result := r.db.Save(oldAssembly)
	return oldAssembly, result.Error
}

func (r *Repository) FormAssembly(formAssembly ds.AssemblyForm) (ds.Assembly, error) {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, formAssembly.Factory_id); result.Error != nil {
		return ds.Assembly{}, result.Error
	}

	if oldAssembly.ID == 0 {
		return ds.Assembly{}, AssemblyNotFound
	}

	oldAssembly.Status = utils.ExistsString
	oldAssembly.DateStartOfProcessing = time.Now()

	result := r.db.Save(oldAssembly)
	return oldAssembly, result.Error
}

func (r *Repository) CompleteAssembly(formAssembly ds.AssemblyForm) (ds.Assembly, error) {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, formAssembly.Factory_id); result.Error != nil {
		return ds.Assembly{}, result.Error
	}

	if oldAssembly.ID == 0 {
		return ds.Assembly{}, AssemblyNotFound
	}

	oldAssembly.Status = utils.Сompleted
	oldAssembly.DateEnd = time.Now()

	result := r.db.Save(oldAssembly)
	return oldAssembly, result.Error
}

func (r *Repository) RejectAssembly(formAssembly ds.AssemblyForm) (ds.Assembly, error) {
	oldAssembly := ds.Assembly{}
	if result := r.db.First(&oldAssembly, formAssembly.Factory_id); result.Error != nil {
		return ds.Assembly{}, result.Error
	}

	if oldAssembly.ID == 0 {
		return ds.Assembly{}, AssemblyNotFound
	}

	oldAssembly.Status = utils.Rejected
	oldAssembly.DateEnd = time.Now()

	result := r.db.Save(oldAssembly)
	return oldAssembly, result.Error
}
