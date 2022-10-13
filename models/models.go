package models

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/shlason/kaigon/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

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

	rd := redis.NewClient(&redis.Options{
		Addr:     configs.Database.Redis.Address,
		Password: configs.Database.Redis.Password,
		DB:       configs.Database.Redis.DB,
	})

	d.AutoMigrate(&Account{})
	db = d
	rdb = rd
}
