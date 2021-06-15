package telepathy

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type WebsocketTelepathy struct {
	ws        *websocket.Conn
	config    *configuration.FreezeConfig
	handlers  []func(string)
	isRunning bool
	stopChan  *chan bool
}

func CreateWebsocketTelepathy(_config *configuration.FreezeConfig) *WebsocketTelepathy {
	_qChan := make(chan bool)
	return &WebsocketTelepathy{isRunning: false, stopChan: &_qChan, config: _config}
}

func listenMsg(tele *WebsocketTelepathy) {
	for {
		_, message, err := tele.ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

func (tele *WebsocketTelepathy) Start() (bool, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:6901"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return false, err
	}

	tele.ws = ws

	go listenMsg(tele)

	return true, nil
}

func (c *WebsocketTelepathy) Stop() (bool, error) {
	err := c.ws.Close()
	*c.stopChan <- true
	close(*c.stopChan)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *WebsocketTelepathy) Send(s string) (bool, error) {
	err := c.ws.WriteMessage(websocket.TextMessage, []byte(s))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *WebsocketTelepathy) RegisterHandler(handler func(string)) {
	c.handlers = append(c.handlers, handler)
}
