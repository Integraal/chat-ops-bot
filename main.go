package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var cal chan string
var users []user
var configuration config

type user struct {
	TelegramId   int
	JiraUsername string
	IcsLink      string
}
type config struct {
	Users []user
}

func init() {
	conf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	json.Unmarshal(conf, &configuration)
}
func main() {

}
