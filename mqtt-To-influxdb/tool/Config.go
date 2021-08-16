package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppName string `json:"app_name"`
	AppMode string `json:"app_mode"`
	AppHost string `json:"app_host"`
	AppPort string `json:"app_port"`

	Database    DatabaseConfig `json:"database"`
	Mqtt        MqttConfig     `json:"mqttconfig"`
	RedisConfig RedisConfig    `json:"redis_config"`
	Infludb     InfludbConfig  `json:"infludb"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
	Charset  string `json:"charset"`
	ShowSql  bool   `json:"show_sql"`
}
type InfludbConfig struct {
	Infurl string `json:"infurl"`
	Token  string `json:"token"`
}

//
type MqttConfig struct {
	Mqttip    string `json:"mqttip"`
	Mqttport  int    `json:"mqttport"`
	MqttUname string `json:"mqttUname"`
	MqttPwd   string `json:"mqttPwd"`
	Clientid  string `json:"clientid"`
	Bugger    bool   `json:"bugger"`
}

//Redis属性定义
type RedisConfig struct {
	Addr     string `json:"addr"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

var _cfg *Config = nil

func GetConfig() *Config {
	return _cfg
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&_cfg); err != nil {
		return nil, err
	}
	return _cfg, nil
}
