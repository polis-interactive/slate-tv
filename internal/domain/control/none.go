package control

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type keyboardInputTracker struct {
	sink chan string
	id   int
}

type keyboardReader struct {
	sinks []keyboardInputTracker
}

var kr = (*keyboardReader)(nil)

func getKeyboardReader(id int) chan string {
	if kr == nil {
		kr = &keyboardReader{
			sinks: make([]keyboardInputTracker, 0),
		}
		go runReader()
	}
	tracker := keyboardInputTracker{
		sink: make(chan string),
		id:   id,
	}
	kr.sinks = append(kr.sinks, tracker)
	return tracker.sink
}

func runReader() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if len(kr.sinks) == 0 {
			return
		}
		for _, kt := range kr.sinks {
			kt.sink <- text
		}
	}
}

type noneController struct {
	*Controller
	reader chan string
}

type noneConfig interface{}

var _ controllerImpl = (*noneController)(nil)

func newNoneController(c *Controller, cfg noneConfig) (*noneController, error) {
	r := getKeyboardReader(0)
	n := &noneController{
		Controller: c,
		reader:     r,
	}
	return n, nil
}

func (n *noneController) RunMainLoop() {

	log.Println("noneControl, Main Loop: running")

	for {
		err := n.tryRunMainLoop()
		if err != nil {
			log.Println(fmt.Sprintf("noneControl, main loop: error %v", err))
		}
		select {
		case _, ok := <-n.shutdowns:
			if !ok {
				goto CloseNoneLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("noneControl, Main Loop: retrying")
		}
	}
CloseNoneLoop:
	log.Println("noneControl, Main Loop: closed")
	n.wg.Done()
}

func (n *noneController) tryRunMainLoop() (err error) {
	lastValue := 0.1
	for {
		select {
		case _, ok := <-n.shutdowns:
			if !ok {
				goto CleanUp
			}
		case _ = <-n.reader:
			if lastValue == 0.1 {
				lastValue = 0.5
			} else if lastValue == 0.5 {
				lastValue = 0.8
			} else {
				lastValue = 0.1
			}
			log.Println(lastValue)
			n.SetInputValue(0, lastValue)
		}
	}

CleanUp:
	return err
}
