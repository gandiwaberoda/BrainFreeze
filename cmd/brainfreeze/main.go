package main

import (
	"fmt"
	"log"
	"sync"
	_ "time"

	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
	"harianugrah.com/brainfreeze/pkg/telepathy"
)

func main() {
	globalWaitGroup := sync.WaitGroup{}

	config, err := configuration.LoadStartupConfig()
	if err != nil {
		log.Fatalln("Gagal meload config", err)
	}

	// Local State
	state := state.CreateStateAccess(&config)
	state.StartWatcher(&config)
	defer state.StopWatcher()

	// Telepathy
	globalWaitGroup.Add(1)
	// telepathyChannel := telepathy.CreateWebsocketTelepathy(&config)
	telepathyChannel := telepathy.CreateConsoleTelepathy()
	_, errTelepathy := telepathyChannel.Start()
	if errTelepathy != nil {
		log.Fatalln(errTelepathy.Error())
	}
	defer telepathyChannel.Stop()
	telepathyChannel.RegisterHandler(func(s string) {
		fmt.Println("handle", s)
		intercom, err := models.ParseIntercom(s)
		if err != nil {
			fmt.Println("Bukan intercom", err)
		} else {
			fmt.Println("INTERCOM", intercom)
		}
	})

	// Telemetry
	globalWaitGroup.Add(1)
	telemetry := diagnostic.CreateNewTelemetry(telepathyChannel, &config, state)
	telemetry.Start()
	defer telemetry.Stop()

	// for i := 0; true; i++ {
	// 	x := models.Transform{WorldXcm: models.Centimeter(i)}
	// 	state.UpdateMyTransform(x)
	// 	fmt.Println("what")
	// 	telepathyChannel.Send(fmt.Sprint("Apalah", i))
	// 	// time.Sleep(time.Second * 1)
	// }

	globalWaitGroup.Wait()
}
