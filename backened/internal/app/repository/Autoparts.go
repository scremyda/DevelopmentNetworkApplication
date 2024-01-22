package repository

import (
	"backened/internal/app/ds"
	"strconv"
	"strings"
)

func (r *Repository) Searchautopart(search string) (*[]ds.Autopart, error) {
	var autoparts []ds.Autopart

	if search == "" {
		if err := r.db.Where("status = ?", "true").Find(&autoparts).Error; err != nil {
			return nil, err
		}
		return &autoparts, nil
	}

	if err := r.db.Where("LOWER(name) LIKE ? AND status = true", "%"+strings.ToLower(search)+"%").Find(&autoparts).Error; err != nil {
		return nil, err
	}

	return &autoparts, nil
}

func (r *Repository) AutopartById(id string) (*ds.Autopart, error) {
	var autoparts ds.Autopart
	intId, _ := strconv.Atoi(id)
	r.db.Find(&autoparts, intId)
	return &autoparts, nil
}

func (r *Repository) Deleteautopart(id string) {
	query := "UPDATE autoparts SET Status = 'false' WHERE id = $1"
	r.db.Exec(query, id)
}
