package conf

import "testing"

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig()
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetDBCfg(t *testing.T) {
	conf, err := LoadConfig()
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = GetDBCfg(conf)
	if err != nil {
		t.Error(err.Error())
		return
	}

	return
}
