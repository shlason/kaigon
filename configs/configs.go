package configs

import (
	"fmt"
	"os"

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

type database struct {
	Protocol string
	Name     string
	Mysql    dbInfo
}

var Server = server{}
var Database = database{}

func init() {
	cfg, err := ini.Load("configs.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
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
			Dialect:  cfg.Section("database.mysql").Key("dialect").String(),
			Username: cfg.Section("database.mysql").Key("username").String(),
			Password: cfg.Section("database.mysql").Key("password").String(),
			Options:  cfg.Section("database.mysql").Key("options").String(),
		},
	}
}
