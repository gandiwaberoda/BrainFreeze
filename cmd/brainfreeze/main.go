package main

import (
	"fmt"
	"time"

	"harianugrah.com/brainfreeze/pkg/telepathy"
)

func main() {
	telepathyChannel := telepathy.CreateConsoleTelepathy()
	telepathyChannel.Start()
	defer telepathyChannel.Stop()

	telepathyChannel.RegisterHandler(func(s string) {
		fmt.Println("handle", s)
	})

	go func() {
		<-time.After(time.Second * 4)
		fmt.Println("???")
		telepathyChannel.Stop()
	}()

	for i := 0; i < 30; i++ {
		telepathyChannel.Send(fmt.Sprint("Apalah", i))
		time.Sleep(time.Second * 1)
	}

	time.Sleep(time.Second * 10)
}
