package models

import (
	"fmt"

	"github.com/shlason/kaigon/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s)/%s?%s",
		configs.Database.Mysql.Username,
		configs.Database.Mysql.Password,
		configs.Database.Protocol,
		configs.Database.Mysql.Address,
		configs.Database.Name,
		configs.Database.Mysql.Options,
	)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	d.AutoMigrate(&Account{})
	db = d
}
