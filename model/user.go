package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm: "type: varchar(20)"`
	Telephone string `gorm: "type: varchar(110); not null; unique"`
	Password string `gorm: "type: size(255); not null"`
}
