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
	"reflect"
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
					log.Debug(string(line))
					var event Event
					err = json.Unmarshal(line, &event)
					if err != nil {
						log.Debug(err)
						return
					}
					log.Debug(event)
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

							for _,task := range(status.TaskManagerStatus) {
								sTask := fmt.Sprintf("%d. %s : %s", task.Id, task.Name, task.Status)
								jsonOutputTextArea := jsDoc.Call("createElement", "li")
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output text area")
									return
								}
								jsonOutputTextArea.Set("innerHTML", sTask)
								jsDoc.Call("getElementById", "taskmanagerstatus").Call("appendChild", jsonOutputTextArea)
							}


							//DisplayKey := []string{"status.", "DaemonRunning", "TotalUptimePercentage"}

							values := reflect.ValueOf(&status).Elem()

							for key := 0; key < values.NumField(); key++ {
								name := values.Type().Field(key).Name
								value := values.Field(key).Interface()

								if (name == "TaskManagerStatus") || (name == "TotalUptimePercentage"){
									continue
								}
								log.Debug(values.Type().Field(key).Name)
								jsonOutputTextArea := jsDoc.Call("getElementById", name)
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output text area")
									return
								}
								//sValue := fmt.Sprintf("%s:%s", name, value)

								jsonOutputTextArea.Set("innerHTML", value)
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "percentage")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return
							}
							jsonOutputTextArea.Set("value", status.TotalUptimePercentage.Percentage)

							jsonOutputTextArea = jsDoc.Call("getElementById", "secondsfrominception")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return
							}
							jsonOutputTextArea.Set("innerHTML", status.TotalUptimePercentage.SecondsFromInception)
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


							values := reflect.ValueOf(&bcnBalance).Elem()

							for key := 0; key < values.NumField(); key++ {
								name := values.Type().Field(key).Name
								value := values.Field(key).Interface()

								if (name == "Owed") || (name == "Owe") || (name == "Id"){
									continue
								}
								log.Debug(values.Type().Field(key).Name)
								jsonOutputTextArea := jsDoc.Call("getElementById", name)
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output text area")
									return
								}

								jsonOutputTextArea.Set("innerHTML", value)
							}


						}

						case "Peers": {
							log.Debug("Peers Hit")

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "Peers")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area")
								return
							}
							sValue := fmt.Sprintf("%s", val)
							jsonOutputTextArea.Set("innerHTML", sValue)


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

							values := reflect.ValueOf(&settings).Elem()

							for key := 0; key < values.NumField(); key++ {
								name := values.Type().Field(key).Name
								value := values.Field(key).Interface()
								log.Debug(name)
								if (name == "NodeIndex") || (name == "DeviceID") || (name == "PublicKey") || (name == "IsDNSEligible") || (name == "DesktopApplicationNotification") || (name == "DesktopApplicationAutoStart") || (name == "DNS") {
									continue
								}

								log.Debug(values.Type().Field(key).Name)
								jsonOutputTextArea := jsDoc.Call("getElementById", name)
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output text area")
									return
								}
								// sValue := fmt.Sprintf("%s:%s", name, value)

								jsonOutputTextArea.Set("innerHTML", value)
							}
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
