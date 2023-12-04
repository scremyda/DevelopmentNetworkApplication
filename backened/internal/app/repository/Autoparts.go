package repository

import (
	"backened/internal/app/ds"
	"strconv"
	"strings"
)

func (r *Repository) Searchautopart(search string) (*[]ds.Autopart, error) {
	var autoparts []ds.Autopart
	if search == "" {
		r.db.Where("Status = ?", "Available").Find(&autoparts)
		return &autoparts, nil
	}
	r.db.Find(&autoparts)

	var filteredautoparts []ds.Autopart
	for _, autopart := range autoparts {
		if strings.Contains(strings.ToLower(autopart.Name), strings.ToLower(search)) {
			filteredautoparts = append(filteredautoparts, autopart)
		}
	}

	return &filteredautoparts, nil
}

func (r *Repository) AutopartById(id string) (*ds.Autopart, error) {
	var autoparts ds.Autopart
	intId, _ := strconv.Atoi(id)
	r.db.Find(&autoparts, intId)
	return &autoparts, nil
}

func (r *Repository) Deleteautopart(id string) {
	query := "UPDATE autoparts SET Status = 'Deleted' WHERE id = $1"
	r.db.Exec(query, id)
}
