package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerName string   `yaml:"ServerName"`
	Port       int      `yaml:"Port"`
	ServerInfo []Server `yaml:"Servers"`
}

type Server struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	Ip   string `yaml:"ip"`
}

// Reads Config file to add to a sruct
func ReadConfig(path string) (*Config, error) {
	// Read file from specific path
	fmt.Println("Reading Config file... ", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file... %v", err)

	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml... %v", err)

	}

	fmt.Println(config.ServerName)

	return &config, nil
}

func main() {
	fmt.Println("Hello world...")

	// TODO: Add Flags

	// TODO: Require path to be a valid YAML file

	// TODO: read config file
	cfg, err := ReadConfig("test.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(*cfg)
}
