package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	data, err := os.ReadFile("configs/database.yaml")
	if err != nil {
		return nil, err
	}

	var DBconfig *DatabaseConfig

	err = yaml.Unmarshal(data, &DBconfig)
	if err != nil {
		return nil, err
	}

	log.Println("Database config readed")

	return DBconfig, nil
}

type HTTPConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Secret string `yaml:"secret"`
}

func LoadHTTPConfig() (*HTTPConfig, error) {
	data, err := os.ReadFile("configs/HTTPserver.yaml")
	if err != nil {
		return nil, err
	}

	var HTTPconfig *HTTPConfig

	err = yaml.Unmarshal(data, &HTTPconfig)
	if err != nil {
		return nil, err
	}

	log.Println("Database config readed")

	return HTTPconfig, nil
}
