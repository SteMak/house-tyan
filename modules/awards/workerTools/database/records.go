package database

import (
	"github.com/jinzhu/gorm"
)

type records struct {
	db *gorm.DB
}

func (r *records) New(item *Record) error {
	return r.db.Create(item).Error
}

func (r *records) Record(id string) (*Record, error) {
	result := &Record{}

	if err := r.db.Where("embed_id = ?", id).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *records) Delete(id string) error {
	return r.db.Where("embed_id = ?", id).Delete(&Record{}).Error
}
