package control

import (
	"fmt"
	"log"
	"time"
)

type noneController struct {
	*Controller
}

type noneConfig interface{}

var _ controllerImpl = (*noneController)(nil)

func newNoneController(c *Controller, cfg noneConfig) (*noneController, error) {
	n := &noneController{
		Controller: c,
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
	lastValue := 0
	ticker := time.NewTicker(20.0 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-n.shutdowns:
			if !ok {
				goto CleanUp
			}
		case _ = <-ticker.C:
			lastValue = (lastValue + 1) % 5
			n.SetInputValue(0, float64(lastValue)/5.0)
		}
	}

CleanUp:
	return err
}
