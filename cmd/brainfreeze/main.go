package main

import (
	"fmt"
	"log"
	"sync"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/telepathy"
)

func main() {
	globalWaitGroup := sync.WaitGroup{}

	config, err := configuration.LoadStartupConfig()
	if err != nil {
		log.Fatalln("Gagal meload config", err)
	}

	// Telepathy
	globalWaitGroup.Add(1)
	telepathyChannel := telepathy.CreateWebsocketTelepathy(&config)
	telepathyChannel.Start()
	defer telepathyChannel.Stop()
	telepathyChannel.RegisterHandler(func(s string) {
		fmt.Println("handle", s)
	})

	globalWaitGroup.Wait()

	// for i := 0; i < 30; i++ {
	// 	telepathyChannel.Send(fmt.Sprint("Apalah", i))
	// 	time.Sleep(time.Second * 1)
	// }
}
