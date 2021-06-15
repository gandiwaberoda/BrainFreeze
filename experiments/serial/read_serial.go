package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/cu.usbmodem14401", Baud: 115200}
	ser, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(ser)
	for scanner.Scan() {
		str := scanner.Text()

		if str[0] != '*' || str[len(str)-1] != '#' {
			fmt.Println("Bad format", scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
	}
	if scanner.Err() != nil {
		log.Fatal(err)
	}
}
