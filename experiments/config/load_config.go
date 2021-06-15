package main

import (
	"fmt"
	"log"

	"harianugrah.com/brainfreeze/pkg/models/local_config"
)

func main() {
	config, err := localconfig.LoadFreezeConfig()
	// _ = config

	if err != nil {
		log.Fatalf("Failed to load local configuration: " + err.Error())
	}

	fmt.Println(config)
}
