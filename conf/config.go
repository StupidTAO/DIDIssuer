package conf

import (
	config "github.com/Unknwon/goconfig"
)

const (
	WORKDIR = "/Users/oker/go/src/github.com/DIDIssuer/"
	CONFIGFILE_PATH = "conf/conf.ini"
	SECTION_DB = "mysql"
	USER = "user"
	PASSWORD = "password"
	URL = "url"
)


type DBConfig struct {
	User string
	Password string
	Url string
}

func LoadConfig() (*config.ConfigFile, error) {
	cfg, err := config.LoadConfigFile(WORKDIR+CONFIGFILE_PATH)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetDBCfg(cfg *config.ConfigFile) (DBConfig, error) {
	var dbCfg DBConfig
	user, err := cfg.GetValue(SECTION_DB, USER)
	if err != nil {
		return DBConfig{}, err
	}
	password, err := cfg.GetValue(SECTION_DB, PASSWORD)
	if err != nil {
		return DBConfig{}, nil
	}
	url, err := cfg.GetValue(SECTION_DB, URL)
	if err != nil {
		return DBConfig{}, nil
	}

	dbCfg.User = user
	dbCfg.Password = password
	dbCfg.Url = url

	return dbCfg, nil
}
