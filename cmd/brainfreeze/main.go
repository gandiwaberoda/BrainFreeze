package main

import (
	"fmt"
	"log"
	"sync"
	_ "time"

	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/gut"
	"harianugrah.com/brainfreeze/internal/migraine"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
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

	// Gut
	// gut := gut.CreateGutSerial()
	gut := gut.CreateGutConsole()
	globalWaitGroup.Add(1)
	gut.RegisterHandler(func(s string) {
		gtb, err := gutmodel.ParseGutToBrain(s)
		if err != nil {
			log.Println("wrong gtb", err)
			return
		}
		state.UpdateGutToBrain(gtb)
	})
	_, errGut := gut.Start()
	if errGut != nil {
		log.Panicln("Gut not yet opened:", errGut.Error())
	}
	defer gut.Stop()

	// Artificial Intellegence
	migraine := migraine.CreateMigraine(&config, gut, state)
	migraine.Start()
	defer migraine.Stop()

	// Telepathy
	globalWaitGroup.Add(1)
	telepathyChannel := telepathy.CreateWebsocketTelepathy(&config)
	// telepathyChannel := telepathy.CreateConsoleTelepathy()
	telepathyChannel.RegisterHandler(func(s string) {
		// fmt.Println("handle", s)
		intercom, err := models.ParseIntercom(s)
		if err != nil {
			fmt.Println("Bukan intercom", err)
			return
		}

		if intercom.Kind == models.COMMAND {
			// Bawa ke migraine
			migraine.AddCommand(intercom)
		}
	})
	_, errTelepathy := telepathyChannel.Start()
	if errTelepathy != nil {
		log.Fatalln(errTelepathy.Error())
	}
	defer telepathyChannel.Stop()

	// Telemetry
	globalWaitGroup.Add(1)
	telemetry := diagnostic.CreateNewTelemetry(telepathyChannel, &config, state)
	telemetry.Start()
	defer telemetry.Stop()

	// Wanda Vision
	globalWaitGroup.Add(1)
	topCamera := acquisition.CreateTopCameraAcquisition(&config, true, true, true)
	topCamera.Start()
	defer topCamera.Stop()

	// Stream Out
	streamout := diagnostic.CreateNewStreamOutDiagnostic(topCamera, &config)
	streamout.StartTopCameraOutput()
	streamout.Start()

	globalWaitGroup.Wait()
}
