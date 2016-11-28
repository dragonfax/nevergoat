package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"

	yaml "gopkg.in/yaml.v2"
)

type SettingsFile struct {
	Token     string `yaml:"token"`
	Notestore string `yaml:"notestore"`
}

var Settings SettingsFile

const SettingsFileName = ".nevergote/settings.yaml"

func DefaultSettingsFileName() string {
	u, err := user.Current()
	if err != nil {
		log.Panicf("Failed to get current user: %v", err)
	}
	return u.HomeDir + "/" + SettingsFileName
}

func readSettings() {

	settingsFileName := DefaultSettingsFileName()
	if _, err := os.Stat(settingsFileName); os.IsNotExist(err) {
		log.Fatalf("Settings file '%s' not found.", settingsFileName)
	}

	content, err := ioutil.ReadFile(settingsFileName)
	if err != nil {
		log.Fatalf("Error reading settings file '%s': %v", settingsFileName, err)
	}

	err = yaml.Unmarshal(content, &Settings)
	if err != nil {
		log.Fatalf("Could not parse settings file '%s': %v", settingsFileName, err)
	}

	if Settings.Notestore == "" {
		log.Fatalf("Settings file missing `notestore` field")
	}

	if Settings.Token == "" {
		log.Fatalf("Settings file missing 'token' field")
	}
}
