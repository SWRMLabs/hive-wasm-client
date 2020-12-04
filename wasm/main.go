package main

import (
	"bufio"
	"encoding/json"
	logger "github.com/ipfs/go-log/v2"
	"net/http"
	"syscall/js"
	"time"
	//"io/ioutil"
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

					var event Event
					err = json.Unmarshal(line, &event)
					if err != nil {
						log.Debug(err)
						return
					}
					//log.Debugf("%+v",event.Result.Val)

					var out Out
					err = json.Unmarshal([]byte(event.Result.Val), &out)
					if err != nil {
						log.Debug(err)
						return
					}
					log.Debugf("%+v %+v",event.Result.Topic,out.Data)
					if event.Result.Topic == "Settlement" {
						val, err := json.Marshal(out.Data)
						if err != nil {
							log.Debug(err)
							return
						}
						var settlement Settlement
						err = json.Unmarshal(val, &settlement)
						if err != nil {
							log.Debug(err)
							return
						}
						log.Debug(settlement.Cycle)
					}

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
