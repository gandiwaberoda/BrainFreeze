package telepathy

import (
	"errors"
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

		for _, handler := range tele.handlers {
			handler(string(message))
		}
	}
}

func (tele *WebsocketTelepathy) Start() (bool, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:6901"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return false, errors.New("failed connecting to websocket: " + err.Error())
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
	if !c.isRunning {
		return false, errors.New("websocket is not running")
	}

	err := c.ws.WriteMessage(websocket.TextMessage, []byte(s))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *WebsocketTelepathy) RegisterHandler(handler func(string)) {
	c.handlers = append(c.handlers, handler)
}
