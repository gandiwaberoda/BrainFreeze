package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/gut"
	"harianugrah.com/brainfreeze/internal/migraine"
	"harianugrah.com/brainfreeze/internal/wanda"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
	"harianugrah.com/brainfreeze/pkg/models/state"
	"harianugrah.com/brainfreeze/pkg/telepathy"
)

func main() {
	config, err := configuration.LoadStartupConfig()
	if err != nil {
		log.Fatalln("Gagal meload config", err)
	}

	selfCheck := diagnostic.ConfigValidate(config)
	if selfCheck != nil {
		fmt.Println(selfCheck)
		return
	}
	fmt.Println("Self check finished")

	// Mulai Proses
	globalWaitGroup := sync.WaitGroup{}

	// Local State
	globalWaitGroup.Add(1)
	state := state.CreateStateAccess(&config)
	state.StartWatcher(&config)
	defer state.StopWatcher()

	// Gut
	var gutTalk gut.GutInterface
	if strings.ToUpper(config.Serial.Ports[0]) == "CONSOLE" {
		gutTalk = gut.CreateGutConsole()
	} else if strings.ToUpper(config.Serial.Ports[0]) == "IGNORE" {
		gutTalk = gut.CreateIgnoreConsole()
	} else {
		gutTalk = gut.CreateGutSerial(&config)
	}
	globalWaitGroup.Add(1)
	gutTalk.RegisterHandler(func(s string) {
		gtb, err := gutmodel.ParseGutToBrain(s)
		if err != nil {
			log.Println("wrong gtb", err)
			return
		}
		state.UpdateGutToBrain(gtb)
	})
	_, errGut := gutTalk.Start()
	if errGut != nil {
		log.Panicln("Gut not yet opened:", errGut.Error())
	}
	defer gutTalk.Stop()

	// Artificial Intellegence
	migraine := migraine.CreateMigraine(&config, gutTalk, state)
	migraine.Start()
	defer migraine.Stop()

	// Telepathy
	globalWaitGroup.Add(1)

	var telepathyChannel telepathy.Telepathy
	if strings.ToUpper(config.Telepathy.ChitChatHost[0]) == "CONSOLE" {
		telepathyChannel = telepathy.CreateConsoleTelepathy()
	} else {
		telepathyChannel = telepathy.CreateWebsocketTelepathy(&config)
	}
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

	// Stream Out
	// streamout := diagnostic.CreateNewStreamOutDiagnostic(topCamera, &config)
	// streamout.StartTopCameraOutput()
	// streamout.Start()

	// Wanda Vision
	// Harus dijalankan paling terakhir, kalau mau nampilin Window di Macos karena bersifat blocking
	globalWaitGroup.Add(1)
	vision := wanda.NewWandaVision(&config, state)
	vision.Start()

	globalWaitGroup.Wait()
}
