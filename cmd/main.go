package main

import (
	"fmt"
	conf "redis-go-dispatcher/config"
	"redis-go-dispatcher/server"
)

func main() {
	configPath := "config.yaml"
	loadedConfig, err := conf.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	server.StartServer(loadedConfig)
}
