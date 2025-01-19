package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Current_user_name string
	Db_url            string
}

func Read() Config {
	file, _ := getConfigFilePath()
	bytes, _ := os.ReadFile(file)
	conf := Config{}
	json.Unmarshal(bytes, &conf)
	return conf
}

func (c *Config) SetUser(name string) {
	c.Current_user_name = name
	write(c)
}

func getConfigFilePath() (string, error) {
	fp, _ := os.Getwd()
	fp += "/" + configFileName
	return fp, nil
}

func write(cfg *Config) error {
	file, _ := getConfigFilePath()
	bytes, _ := json.MarshalIndent(*cfg, "", "    ")
	os.WriteFile(file, bytes, 0644)
	return nil
}
