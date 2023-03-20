package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20)"`
	Telephone string `gorm:"type:varchar(30);unique_index"`
	Password  string `gorm:"type:varchar(255)"`
}
