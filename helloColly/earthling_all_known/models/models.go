package models

import (
	"earthling_all_known/dao"

	"github.com/jinzhu/gorm"
)

// License 车牌信息
type License struct {
	gorm.Model
	Identifier  string
	ReleaseDate string
	Rating      string
	Actor       string
	CoverPath   string
}

// CreateALicense 新增一个车牌
func CreateALicense(lic *License) (err error) {
	err = dao.DB.Create(lic).Error
	return
}
