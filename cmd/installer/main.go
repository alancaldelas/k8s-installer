package main

import (
	"fmt"
)

func readConfig(path string) {
	// Read file from specific path
	fmt.Println("Reading Config file... ", path)
}

func main() {
	fmt.Println("Hello world...")

	// TODO: Add Flags

	// TODO: Require path to be a valid YAML file

	// TODO: read config file
	readConfig("Some Path")
}
