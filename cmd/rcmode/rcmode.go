package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/gut"
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

	// selfCheck := diagnostic.ConfigValidate(config)
	// if selfCheck != nil {
	// 	fmt.Println(selfCheck)
	// 	return
	// }
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

		t := models.Transform{
			EncXcm: gtb.AbsX,
			EncYcm: gtb.AbsY,
			EncROT: gtb.Gyro,
		}
		t.InjectWorldTransfromFromEncTransform(&config)
		state.UpdateMyTransform(t)
	})
	_, errGut := gutTalk.Start()
	if errGut != nil {
		log.Panicln("Gut not yet opened:", errGut.Error())
	}
	defer gutTalk.Stop()

	globalWaitGroup.Add(1)
	rcserv := CreateRcWsWs(&config)
	rcserv.RegisterHandler(func(s string) {
		// fmt.Println(s)

		if len(s) < 2 {
			return
		}

		if s[0:2] == "RC" {
			toFwd := s[2:]
			fmt.Println("To Forward: " + toFwd)
			gutTalk.Send(toFwd)
			// gutTalk.Send("*a,0,0,0,0,1#")
		}

	})
	rcserv.Start()

	// Artificial Intellegence
	// migraine := migraine.CreateMigraine(&config, gutTalk, state)
	// migraine.Start()
	// defer migraine.Stop()

	// Telepathy
	globalWaitGroup.Add(1)

	var telepathyChannel telepathy.Telepathy
	if strings.ToUpper(config.Telepathy.ChitChatHost[0]) == "CONSOLE" {
		telepathyChannel = telepathy.CreateConsoleTelepathy()
	} else {
		telepathyChannel = telepathy.CreateWebsocketTelepathy(&config)
	}

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
	// globalWaitGroup.Add(1)
	// vision := wanda.NewWandaVision(&config, state)
	// vision.Start()

	globalWaitGroup.Wait()
}

type RcWs struct {
	server    *http.Server
	config    *configuration.FreezeConfig
	handlers  []func(string)
	clients   map[*websocket.Conn]bool
	isRunning bool
}

func CreateRcWsWs(_config *configuration.FreezeConfig) *RcWs {
	return &RcWs{isRunning: false, config: _config, clients: make(map[*websocket.Conn]bool)}
}

var upgrader = websocket.Upgrader{} // use default options

func (s *RcWs) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	s.clients[c] = true

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", message)
		// err = c.WriteMessage(mt, message)

		for _, v := range s.handlers {
			v(string(message))
		}

		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	delete(s.clients, c)
}

func (s *RcWs) Start() (bool, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.echo)

	s.server = &http.Server{
		Addr:    ":6969",
		Handler: mux,
	}

	go func() {
		log.Fatal(s.server.ListenAndServe())
	}()
	return true, nil
}

func (c *RcWs) Stop() (bool, error) {
	return false, nil
}

func (c *RcWs) Broadcast(s string) (bool, error) {
	for v := range c.clients {
		v.WriteMessage(websocket.TextMessage, []byte(s))
	}
	return true, nil
}

func (c *RcWs) RegisterHandler(handler func(string)) {
	c.handlers = append(c.handlers, handler)
}
