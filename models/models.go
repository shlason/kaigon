package models

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/shlason/kaigon/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mongoDBCollections struct {
	ChatMessages *mongo.Collection
}

var db *gorm.DB
var rdb *redis.Client
var mdb mongoDBCollections = mongoDBCollections{}

var rctx context.Context = context.Background()

func init() {
	// MySQL
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
	d.AutoMigrate(&AccountOauthInfo{})
	d.AutoMigrate(&AccountProfile{})
	d.AutoMigrate(&AccountSetting{})
	d.AutoMigrate(&AccountSettingNotification{})
	d.AutoMigrate(&AccountProfileSocialMedia{})
	d.AutoMigrate(&ChatRoom{})
	d.AutoMigrate(&ChatRoomMember{})

	// MongoDB
	md, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(fmt.Sprintf("%s://%s", configs.Database.MongoDB.OProtocol, configs.Database.MongoDB.Address)),
	)

	if err != nil {
		panic(err)
	}

	if err := md.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	chatMessagesColl := md.Database(configs.Database.Name).Collection(chatMessagesCollectionName)

	// Redis
	rd := redis.NewClient(&redis.Options{
		Addr:     configs.Database.Redis.Address,
		Password: configs.Database.Redis.Password,
		DB:       configs.Database.Redis.DB,
	})

	db = d
	rdb = rd
	mdb = mongoDBCollections{
		ChatMessages: chatMessagesColl,
	}
}
