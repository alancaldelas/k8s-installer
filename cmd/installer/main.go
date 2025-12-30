package main

import (
	"fmt"
	"os"

	"github.com/alancaldelas/k8s-installer/pkg/config"
	"github.com/alancaldelas/k8s-installer/pkg/libvirt"
	"github.com/alancaldelas/k8s-installer/pkg/network"
	"github.com/alancaldelas/k8s-installer/pkg/ssh"
)

func main() {
	// TODO: Add Flags

	// TODO: Require path to be a valid YAML file

	// TODO: read config file
	cfg, err := config.ReadConfig("test.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", *cfg)

	// TODO: Connect to each server
	for _, server := range cfg.ServerInfo {
		fmt.Printf("Checking connectivity to %s (%s:%d)...\n", server.Name, server.Ip, server.Port)
		if err := network.CheckConnectivity(server.Ip, server.Port); err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", server.Name, err)
		} else {
			fmt.Printf("Successfully connected to %s\n", server.Name)

			// Determine user (Server override > Global config)
			user := cfg.User
			if server.User != "" {
				user = server.User
			}

			// TODO: Run a simple command to verify connectivity
			sshClient, err := ssh.GetSSHClient(user, cfg.PrivateKeyPath, server.Ip, server.Port)
			if err != nil {
				fmt.Printf("Failed to create SSH client for %s: %v\n", server.Name, err)
				continue
			}
			defer sshClient.Close()

			if err := ssh.RunCommand(sshClient, "echo Hello"); err != nil {
				fmt.Printf("Failed to run command on %s: %v\n", server.Name, err)
			} else {
				fmt.Printf("Successfully ran command on %s\n", server.Name)
				// TODO: List VMs running on the server
				if err := libvirt.ListDomains(user, cfg.PrivateKeyPath, server.Ip); err != nil {
					fmt.Printf("Failed to list VMs on %s: %v\n", server.Name, err)
				}
			}
		}
	}
}
