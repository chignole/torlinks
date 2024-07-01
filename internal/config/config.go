package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type configuration struct {
	General general `json:"general"`
}

type general struct {
	Source                string   `json:"source"`
	Destination           string   `json:"destination"`
	Data                  []string `json:"data"`
	UpdateDatabaseOnStart string   `json:"updatedatabaseonstart"`
}

// Processing configuration file
func InitConfig() configuration {
	var loadConfig configuration

	log.Println("[INFO] Loading configuration file...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	file := homeDir + "/.config/torlinks/config.json"
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("[ERROR] Error opening configuration file : %v", err)
	}
	defer configFile.Close()
	config, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Panicf("[ERROR] Error processing configuration file : %v", err)
	}
	err = json.Unmarshal(config, &loadConfig)
	if err != nil {
		log.Panicf("[ERROR] Error processing configuration file : %v", err)
	}
	return loadConfig
}
