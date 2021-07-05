package simpserv

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type SimpWs struct {
	server    *http.Server
	config    *configuration.FreezeConfig
	handlers  []func(string)
	clients   map[*websocket.Conn]bool
	isRunning bool
}

func CreateSimpWs(_config *configuration.FreezeConfig) *SimpWs {
	return &SimpWs{isRunning: false, config: _config, clients: make(map[*websocket.Conn]bool)}
}

var upgrader = websocket.Upgrader{} // use default options

func (s *SimpWs) echo(w http.ResponseWriter, r *http.Request) {
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

func (s *SimpWs) Start() (bool, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.echo)

	s.server = &http.Server{
		Addr:    s.config.Simulator.SimpservPort,
		Handler: mux,
	}

	go func() {
		log.Fatal(s.server.ListenAndServe())
	}()
	return true, nil
}

func (c *SimpWs) Stop() (bool, error) {
	return false, nil
}

func (c *SimpWs) Broadcast(s string) (bool, error) {
	for v := range c.clients {
		v.WriteMessage(websocket.TextMessage, []byte(s))
	}
	return true, nil
}

func (c *SimpWs) RegisterHandler(handler func(string)) {
	c.handlers = append(c.handlers, handler)
}
