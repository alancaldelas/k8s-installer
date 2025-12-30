package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"libvirt.org/go/libvirt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type Config struct {
	User           string   `yaml:"User"`
	PrivateKeyPath string   `yaml:"PrivateKeyPath"`
	ServerInfo     []Server `yaml:"Servers"`
}

type Server struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	Ip   string `yaml:"Ip"`
	User string `yaml:"User"`
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

func main() {
	// TODO: Add Flags

	// TODO: Require path to be a valid YAML file

	// TODO: read config file
	cfg, err := ReadConfig("test.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", *cfg)

	// TODO: Connect to each server
	for _, server := range cfg.ServerInfo {
		fmt.Printf("Checking connectivity to %s (%s:%d)...\n", server.Name, server.Ip, server.Port)
		if err := CheckConnectivity(server.Ip, server.Port); err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", server.Name, err)
		} else {
			fmt.Printf("Successfully connected to %s\n", server.Name)
			
			// Determine user (Server override > Global config)
			user := cfg.User
			if server.User != "" {
				user = server.User
			}

			// TODO: Run a simple command to verify connectivity
			sshClient, err := GetSSHClient(user, cfg.PrivateKeyPath, server.Ip, server.Port)
			if err != nil {
				fmt.Printf("Failed to create SSH client for %s: %v\n", server.Name, err)
				continue
			}
			defer sshClient.Close()
			
			if err := RunCommand(sshClient, "echo Hello"); err != nil {
				fmt.Printf("Failed to run command on %s: %v\n", server.Name, err)
			} else {
				fmt.Printf("Successfully ran command on %s\n", server.Name)
				// TODO: List VMs running on the server
				if err := ListDomains(user, cfg.PrivateKeyPath, server.Ip); err != nil {
					fmt.Printf("Failed to list VMs on %s: %v\n", server.Name, err)
				}
			}
		}
	}
}

func GetSSHClient(user, privateKeyPath, ip string, port int) (*ssh.Client, error) {
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return client, nil
}

func RunCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	fmt.Println(string(output))
	return nil
}

func CheckConnectivity(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func ListDomains(user, keyPath, ip string) error {
	// e.g. qemu+ssh://alan@10.0.129.11/system?keyfile=/home/alan/.ssh/id_rsa&no_verify=1
	uri := fmt.Sprintf("qemu+ssh://%s@%s/system?keyfile=%s&no_verify=1", user, ip, keyPath)
	
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return fmt.Errorf("failed to list domains: %v", err)
	}

	fmt.Printf("%d VMs found:\n", len(domains))
	for _, d := range domains {
		name, _ := d.GetName()
		id, _ := d.GetID()
		state, _, _ := d.GetState()
		stateStr := "Unknown"
		switch state {
		case libvirt.DOMAIN_RUNNING:
			stateStr = "Running"
		case libvirt.DOMAIN_SHUTOFF:
			stateStr = "Shutoff"
		case libvirt.DOMAIN_PAUSED:
			stateStr = "Paused"
		}
		fmt.Printf(" - %s (ID: %d, State: %s)\n", name, id, stateStr)
		d.Free()
	}

	return nil
}
