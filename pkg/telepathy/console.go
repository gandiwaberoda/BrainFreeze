package telepathy

import (
	"bufio"
	"fmt"
	"os"

	frerror "harianugrah.com/brainfreeze/pkg/errors"
)

type ConsoleTelepathy struct {
	handlers  []func(string)
	isRunning bool
	stopChan  *chan bool
}

func CreateConsoleTelepathy() *ConsoleTelepathy {
	_qChan := make(chan bool)
	return &ConsoleTelepathy{isRunning: false, stopChan: &_qChan}
}

func (c *ConsoleTelepathy) Start() (bool, error) {
	c.isRunning = true

	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			select {
			case doStop := <-*c.stopChan:
				fmt.Println("doStop", doStop)

				if doStop {
					c.isRunning = false
					fmt.Println("Stopped")
					close(*c.stopChan)
					return
				}
			default:
				fmt.Print("Enter text: ")
				text, _ := reader.ReadString('\n')
				for _, v := range c.handlers {
					v(text)
				}
			}

		}
	}()

	return true, nil
}

func (c *ConsoleTelepathy) Stop() (bool, error) {
	fmt.Println("Panggil stop")
	*c.stopChan <- true
	return true, nil
}

func (c *ConsoleTelepathy) Send(s string) (bool, error) {
	if !c.isRunning {
		return false, &frerror.TelepathyNotRunningError{}
	}

	fmt.Println(s)
	return true, nil
}

func (c *ConsoleTelepathy) RegisterHandler(handler func(string)) {
	c.handlers = append(c.handlers, handler)
}
