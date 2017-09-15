package config

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/integraal/chat-ops-bot/telegram"
	"github.com/integraal/chat-ops-bot/components/user"
)

const (
	confFileName = "config.json"
)

type Config struct {
	Users    []user.User
	Telegram telegram.Config
}

func Read(filename string) *Config {
	var configuration Config
	conf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	json.Unmarshal(conf, &configuration)
	return &configuration
}

func Initialize() *Config {
	return Read(confFileName)
}
