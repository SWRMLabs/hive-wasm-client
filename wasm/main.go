// GOOS=js GOARCH=wasm go build -o  ../assets/hive.wasm

package main

import (
	"bufio"
	"encoding/json"
	logger "github.com/ipfs/go-log/v2"
	"net/http"
	"syscall/js"
	"time"
	"fmt"
	//"io/ioutil"
)

var log = logger.Logger("events")

const (
	EVENTS = "http://localhost:4343/v3/events"
)


func Events() js.Func {

	jsonfunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {


		go func() {
			resp, err := http.Post(EVENTS, "application/json", nil)
			if err != nil {
				log.Error(err.Error())
				return
			}

			reader := bufio.NewReader(resp.Body)

			for {
				line, _, err := reader.ReadLine()
				if err != nil || line == nil {
					log.Debug("Update Complete")
					Events()
				} else {

					var event Event
					err = json.Unmarshal(line, &event)
					if err != nil {
						log.Debug(err)
						return
					}

					var out Out
					err = json.Unmarshal([]byte(event.Result.Val), &out)
					if err != nil {
						log.Debug(err)
						return
					}
					log.Debugf("%+v %+v",event.Result.Topic, out.Data)

					val, err := json.MarshalIndent(out.Data, "", " ")
					if err != nil {
						log.Debug(err)
						log.Debug("Error encountered")
						return
					}

					switch event.Result.Topic {

						case "Status": {
							log.Debug("Status Hit")
							var status Status
							err = json.Unmarshal(val, &status)
							if err != nil {
								log.Debug(err)
								return
							}
							log.Debug(status)

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "status")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return

							}

							sStatus := fmt.Sprintf("%+v", status)

							jsonOutputTextArea.Set("innerHTML", sStatus)

						}


						case "BalanceCycle": {
							log.Debug("BCN Hit")
							var bcnBalance BCNBalance
							err = json.Unmarshal(val, &bcnBalance)
							if err != nil {
								log.Debug(err)
								return
							}
							log.Debug(bcnBalance)

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "bcnbalance")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return

							}

							sBcnBalance := fmt.Sprintf("%+v", bcnBalance)

							jsonOutputTextArea.Set("innerHTML", sBcnBalance)

						}

						case "Settings": {
							log.Debug("Settings Hit")
							var settings Settings

							err = json.Unmarshal(val, &settings)
							if err != nil {
								log.Debug(err)
								return
							}
							log.Debug(settings)

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "settings")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return

							}

							sSettings := fmt.Sprintf("%+v", settings)
							jsonOutputTextArea.Set("innerHTML", sSettings)

						}
						default:{
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "state")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return

							}
							jsonOutputTextArea.Set("innerHTML", "Process Halted")
						}
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
