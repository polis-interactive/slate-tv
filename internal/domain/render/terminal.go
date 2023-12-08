package render

import (
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
)

type terminalRender struct {
	*baseRender
	pb []util.Color
}

var _ render = (*terminalRender)(nil)

func newTerminalRender(base *baseRender) *terminalRender {

	log.Println("terminalRender, newTerminalRender: creating")

	r := &terminalRender{
		baseRender: base,
	}

	log.Println("terminalRender, newTerminalRender: created")

	return r
}

func (r *terminalRender) runMainLoop() {
	r.pb = make([]util.Color, r.ledCount)
	for {
		err := r.runRenderLoop()
		if err != nil {
			log.Println(fmt.Sprintf("terminal, Main Loop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-r.shutdowns:
			if !ok {
				goto CloseTerminalLoop
			}
		}
	}

CloseTerminalLoop:
	log.Println("terminalRender runMainLoop, Main Loop: closed")
	r.wg.Done()
}

func (r *terminalRender) runRender() error {

	err := r.bus.CopyLightsToColorBuffer(r.pb)
	if err != nil {
		return err
	}

	//outputString := "START(\n"
	//
	//log.Println(r.pb)
	//
	//outputString += ")END"

	// log.Println(outputString)

	return nil
}
