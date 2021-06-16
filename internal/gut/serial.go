package gut

import (
	"bufio"
	"errors"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

type GutSerial struct {
	Gut
	Port *serial.Port
}

func CreateGutSerial() *GutSerial {
	return &GutSerial{}
}

func worker(gut *GutSerial) {
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

func (g *GutSerial) Start() (bool, error) {
	c := &serial.Config{Name: "/dev/cu.usbmodem14401", Baud: 115200}

	ser, err := serial.OpenPort(c)
	if err != nil {
		return false, err
	}

	g.Port = ser
	go worker(g)

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

	_, err := g.Port.Write([]byte(msg))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (g *GutSerial) RegisterHandler(handler func(string)) {
	g.handlers = append(g.handlers, handler)
}
