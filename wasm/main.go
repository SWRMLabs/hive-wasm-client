package main

import (
	"bufio"
	"encoding/json"
	logger "github.com/ipfs/go-log/v2"
	"net/http"
	"syscall/js"
	"time"
	"github.com/StreamSpace/hive-wasm-client/types"
)

var log = logger.Logger("events")

const (
	EVENTS = "http://localhost:4343/v3/events"
)


func Events() js.Func {

	jsonfunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		jsDoc := js.Global().Get("document")
		if !jsDoc.Truthy() {
			log.Debug("Unable to get document object")
			return nil
		}

		jsonOutputTextArea := jsDoc.Call("getElementById", "count")
		if !jsonOutputTextArea.Truthy() {
			log.Debug("Unable to get output text area")
			return nil

		}

		go func() {
			resp, err := http.Post(EVENTS, "application/json", nil)
			if err != nil {
				log.Error(err.Error())
				return
			}
			reader := bufio.NewReader(resp.Body)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					log.Debug("Update Complete")
					return
				} else {
					//log.Debug("I just received the message %s", string(line))

					data := make(map[string]map[string]string)
					json.Unmarshal(line, &data)
					log.Debug(data)
					udata := data["result"]["val"]
					log.Debug(udata)

					time.Sleep(1 * time.Second)
				}

			}
		}()
		return nil
	})
	return jsonfunc
}

func main() {
	logger.SetLogLevel("*", "Debug")
	js.Global().Set("Events", Events())
	<-make(chan bool)
}
