package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type server struct {
	Protocol string
	Host     string
	Port     string
}

type dbInfo struct {
	Dialect  string
	Username string
	Password string
	Address  string
	Options  string
}

type redisDbInfo struct {
	dbInfo
	DB int
}

type database struct {
	Protocol string
	Name     string
	Mysql    dbInfo
	Redis    redisDbInfo
}

type smtp struct {
	Sender   string
	Password string
	Host     string
	Port     string
}

var Server = server{}
var Database = database{}
var Smtp = smtp{}

func init() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("Fail to get user home dir: %s", err)
	}

	cfg, err := ini.Load(filepath.Join(homeDir, os.Getenv("CONFIG_FILE_PATH_BASE_ON_HOME")))
	fmt.Println(filepath.Join(homeDir, os.Getenv("CONFIG_FILE_PATH_BASE_ON_HOME")))
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	redisDb, err := cfg.Section("database.redis").Key("DB").Int()

	if err != nil {
		fmt.Printf("Fail to read redis key 'db': %v", err)
		os.Exit(1)
	}

	Server = server{
		Protocol: cfg.Section("server").Key("protocol").String(),
		Host:     cfg.Section("server").Key("host").String(),
		Port:     cfg.Section("server").Key("port").String(),
	}
	Database = database{
		Protocol: cfg.Section("database").Key("protocol").String(),
		Name:     cfg.Section("database").Key("name").String(),
		Mysql: dbInfo{
			Address:  cfg.Section("database.mysql").Key("address").String(),
			Dialect:  cfg.Section("database.mysql").Key("dialect").String(),
			Username: cfg.Section("database.mysql").Key("username").String(),
			Password: cfg.Section("database.mysql").Key("password").String(),
			Options:  cfg.Section("database.mysql").Key("options").String(),
		},
		Redis: redisDbInfo{
			dbInfo: dbInfo{
				Address:  cfg.Section("database.redis").Key("address").String(),
				Password: cfg.Section("database.redis").Key("password").String(),
			},
			DB: redisDb,
		},
	}
	Smtp = smtp{
		Sender:   cfg.Section("smtp").Key("sender").String(),
		Password: cfg.Section("smtp").Key("google_app_password").String(),
		Host:     cfg.Section("smtp").Key("host").String(),
		Port:     cfg.Section("smtp").Key("port").String(),
	}
}
