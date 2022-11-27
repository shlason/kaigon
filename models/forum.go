package models

import "gorm.io/gorm"

type Forum struct {
	gorm.Model
	Name        string
	Icon        string
	Banner      string
	Rule        string
	Description string
}
