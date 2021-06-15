package main

import (
	"fmt"
	"log"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

func main() {
	config, err := configuration.LoadStartupConfig()
	if err != nil {
		log.Fatalf("Failed to load local configuration: " + err.Error())
	}

	fmt.Println(config)
}
