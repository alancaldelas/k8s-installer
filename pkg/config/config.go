package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User           string    `yaml:"User"`
	PrivateKeyPath string    `yaml:"PrivateKeyPath"`
	ServerInfo     []Server  `yaml:"Servers"`
	ClusterInfo    []Cluster `yaml:"Cluster"`
}

type Server struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	Ip   string `yaml:"Ip"`
	User string `yaml:"User"`
}

type Cluster struct {
	Name    string    `yaml:"name"`
	Network []Network `yaml:"network"`
}

type Network struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	DNS     string `yaml:"dns"`
	Gateway string `yaml:"gateway"`
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

	return &config, nil
}
