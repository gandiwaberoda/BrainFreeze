package gut

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GutSerial struct {
	Gut
	Curstate *state.StateAccess
	Port     *serial.Port
	conf     *configuration.FreezeConfig
	toSend   string
}

func CreateGutSerial(conf *configuration.FreezeConfig, curstate *state.StateAccess) *GutSerial {
	return &GutSerial{
		conf:     conf,
		Curstate: curstate,
	}
}

func workerReader(gut *GutSerial) {
	scanner := bufio.NewScanner(gut.Port)
	for scanner.Scan() {
		str := scanner.Text()

		if len(str) < 2 {
			// Kurang dari dua huruf, salah format
			fmt.Println("Really bad format", str)
		}

		if str[0] != '*' || str[len(str)-1] != '#' {
			fmt.Println("Bad format", str)
		} else {
			for _, handler := range gut.handlers {
				handler(str)
			}
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err().Error())
	}
}

func workerWriter(gut *GutSerial) {
	msDelay := int(time.Second) / gut.conf.Serial.CommandHz
	for {
		<-time.After(time.Duration(msDelay))
		if gut.toSend == "" {
			fmt.Println("EMPTY")
		}

		// fmt.Println("SEND: ", gut.toSend)
		// gut.Send(gut.toSend)

		if _, err := gut.Port.Write([]byte(gut.toSend)); err != nil {
			fmt.Println("FAILED SENDING GUT PORT")
		}
		if !strings.EqualFold(gut.conf.Serial.ArayaPorts[0], "DISABLED") {
			if _, err := gut.ArayaSens.Port.Write([]byte(gut.toSend)); err != nil {
				fmt.Println("FAILED SENDING ARAYA PORT")
			}
		}
		gut.Curstate.UpdateToGutCmd(gut.toSend)
	}
}

func (g *GutSerial) Start() (bool, error) {
	c := &serial.Config{Name: g.conf.Serial.Ports[0], Baud: 9600}

	ser, err := serial.OpenPort(c)
	if err != nil {
		return false, err
	}

	g.Port = ser
	go workerReader(g)
	go workerWriter(g)

	return true, nil
}

func (g *GutSerial) Stop() (bool, error) {
	err := g.Port.Close()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (g *GutSerial) Send(msg string) (bool, error) {
	if g.Port == nil {
		return false, errors.New("port is not yet opened")
	}
	g.toSend = msg
	return true, nil

	// _, err := g.Port.Write([]byte(msg))
	// if err != nil {
	// 	return false, err
	// } else {
	// 	return true, nil
	// }
}

func (g *GutSerial) RegisterHandler(handler func(string)) {
	g.handlers = append(g.handlers, handler)
}
