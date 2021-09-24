package araya

import (
	"bufio"
	"errors"
	"fmt"
	"log"

	"github.com/tarm/serial"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type ArayaSerial struct {
	Curstate *state.StateAccess
	Port     *serial.Port
	conf     *configuration.FreezeConfig
}

func NewArayaSerial(conf *configuration.FreezeConfig, curstate *state.StateAccess) *ArayaSerial {
	return &ArayaSerial{
		conf:     conf,
		Curstate: curstate,
	}
}

func (a *ArayaSerial) onArayaReceived(s string) {
	parsed, err := gutmodel.ParseAraya(s)
	if err != nil {
		fmt.Println("Araya error: ", err)
	}

	a.Curstate.UpdateAraya(parsed)
}

func workerReader(araya *ArayaSerial) {
	scanner := bufio.NewScanner(araya.Port)
	for scanner.Scan() {
		str := scanner.Text()

		if len(str) < 2 {
			// Kurang dari dua huruf, salah format
			fmt.Println("Araya really bad format", str)
		}

		if str[0] != '*' || str[len(str)-1] != '#' {
			fmt.Println("Bad format", str)
		} else {
			// Update araya state
			fmt.Println("Nerima ", str)
			araya.onArayaReceived(str)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err().Error())
	}
}

func (g *ArayaSerial) Start() (bool, error) {
	fmt.Println("Started")

	c := &serial.Config{Name: g.conf.Serial.ArayaPorts[0], Baud: 9600}

	ser, err := serial.OpenPort(c)
	if err != nil {
		return false, err
	}

	g.Port = ser
	go workerReader(g)

	return true, nil
}

func (g *ArayaSerial) Send(msg string) (bool, error) {
	if g.Port == nil {
		return false, errors.New("port araya is not yet opened")
	}
	g.Port.Write([]byte(msg))
	return true, nil
}

func (g *ArayaSerial) Stop() (bool, error) {
	err := g.Port.Close()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
