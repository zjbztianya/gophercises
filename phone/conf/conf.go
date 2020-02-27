package conf

import (
	"github.com/go-ini/ini"
)

var Conf = &Config{}

type Config struct {
	Type      string
	User      string
	Password  string
	Host      string
	DbName    string
	TableName string
}

// Init conf.
func Init() error {
	cfg, err := ini.Load("phone/conf/app.ini")
	if err != nil {
		return err
	}
	return cfg.Section("database").MapTo(Conf)
}
