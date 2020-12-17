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
	"strings"
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
					// log.Debugf("%+v %+v",event.Result.Topic, out.Data)

					val, err := json.MarshalIndent(out.Data, "", " ")
					if err != nil {
						log.Debug(err)
						log.Debug("Error encountered")
						return
					}

					if (event.Result.Topic == "Status") || (event.Result.Topic == "Balance") || (event.Result.Topic == "BalanceCycle") || (event.Result.Topic == "Peers") || (event.Result.Topic == "Settings")  {
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
								log.Debug("Unable to get document object in status")
								return
							}

							for _,task := range(status.TaskManagerStatus) {
								sName := fmt.Sprintf("%s",task.Name)
								jsonOutputTextArea := jsDoc.Call("getElementById", "taskmanagerstatusname")
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output area task manager name")
									return
								}
								//
								jsonOutputTextArea.Set("innerHTML", sName)

								sStatus := fmt.Sprintf("%s",task.Status)
								jsonOutputTextArea = jsDoc.Call("getElementById", "taskmanagerstatusstatus")
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output area task manager status")
									return
								}
								jsonOutputTextArea.Set("innerHTML", sStatus)

								sAdditionalStatus := fmt.Sprintf("%s",task.AdditionalStatus)
								jsonOutputTextArea = jsDoc.Call("getElementById", "taskmanagerstatusAS")
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output area task manager additional status")
									return
								}
								jsonOutputTextArea.Set("innerHTML", sAdditionalStatus)
								// jsDoc.Call("getElementById", "taskmanagerstatus").Call("appendChild", jsonOutputTextArea)
							}


							//DisplayKey := []string{"status.", "DaemonRunning", "TotalUptimePercentage"}

							values := reflect.ValueOf(&status).Elem()

							for key := 0; key < values.NumField(); key++ {
								name := values.Type().Field(key).Name
								value := values.Field(key).Interface()

								if (name == "TaskManagerStatus") || (name == "TotalUptimePercentage") || (name == "SessionStartTime") || (name == "ServerStatus"){
									continue
								}

								log.Debug(values.Type().Field(key).Name)
								jsonOutputTextArea := jsDoc.Call("getElementById", name)
								if !jsonOutputTextArea.Truthy() {
									log.Debug("Unable to get output text area in status keys")
									return
								}
								var sValue string
								if value == true {
								switch name {
								case "LoggedIn":
									sValue = "LoggedIn"

								case "DaemonRunning":
									sValue = "ONLINE"
								}
							} else if value == false{
								switch name {
								case "LoggedIn":
									sValue = "LoggedOut"

								case "DaemonRunning":
									sValue = "OFFLINE"
								}
							}
								//sValue := fmt.Sprintf("%s:%s", name, value)
								jsonOutputTextArea.Set("innerHTML",sValue)

							}

							// jsonOutputTextArea := jsDoc.Call("getElementById", "percentage")
							// if !jsonOutputTextArea.Truthy() {
							// 	log.Debug("Unable to get output text area in percentage")
							// 	return
							// }
							// jsonOutputTextArea.Set("value", status.TotalUptimePercentage.Percentage)

							jsonOutputTextArea := jsDoc.Call("getElementById", "percentageNumber")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area in percentagenumber")
								return
							}
							sFloat := fmt.Sprintf("%.2f", status.TotalUptimePercentage.Percentage)
							sValue := fmt.Sprintf("%s %s", sFloat, "%")
							jsonOutputTextArea.Set("innerHTML", sValue)

							// jsonOutputTextArea = jsDoc.Call("getElementById", "secondsfrominception")
							// if !jsonOutputTextArea.Truthy() {
							// 	log.Debug("Unable to get output text area in secondsfromincepti")
							// 	return
							// }
							// jsonOutputTextArea.Set("innerHTML", status.TotalUptimePercentage.SecondsFromInception)
						}

						case "Balance": {
							log.Debug("Balance Hit")

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object in balance")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "confirmedBalance")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area in balance")
								return
							}
							sFloat := fmt.Sprintf("%s", val)
							for i, value := range(sFloat){
								if (strings.ContainsAny(string(value), ".") && (i+3) <= len(sFloat)){
									sFloat = sFloat[0:i+1] + sFloat[i+1 : i+3]
									break
								}
							}
							sValue := fmt.Sprintf("%s %s", sFloat, "SWRM")
							jsonOutputTextArea.Set("innerHTML", sValue)
						}

						// case "BalanceCycle": {
						// 	log.Debug("BCN Hit")
						// 	var bcnBalance BCNBalance
						// 	err = json.Unmarshal(val, &bcnBalance)
						// 	if err != nil {
						// 		log.Debug(err)
						// 		return
						// 	}
						// 	log.Debug(bcnBalance)
						//
						// 	jsDoc := js.Global().Get("document")
						// 	if !jsDoc.Truthy() {
						// 		log.Debug("Unable to get document object in balance cycle")
						// 		return
						// 	}
						//
						//
						// 	values := reflect.ValueOf(&bcnBalance).Elem()
						//
						// 	for key := 0; key < values.NumField(); key++ {
						// 		name := values.Type().Field(key).Name
						// 		value := values.Field(key).Interface()
						//
						// 		if (name == "Owed") || (name == "Owe") || (name == "Id"){
						// 			continue
						// 		}
						// 		log.Debug(values.Type().Field(key).Name)
						// 		jsonOutputTextArea := jsDoc.Call("getElementById", name)
						// 		if !jsonOutputTextArea.Truthy() {
						// 			log.Debug("Unable to get output text area in balance cycle keys")
						// 			return
						// 		}
						//
						// 		jsonOutputTextArea.Set("innerHTML", value)
						// 	}
						//
						//
						// }

						case "Peers": {
							log.Debug("Peers Hit")

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Debug("Unable to get document object in peers")
								return
							}

							jsonOutputTextArea := jsDoc.Call("getElementById", "PeersData")
							if !jsonOutputTextArea.Truthy() {
								log.Debug("Unable to get output text area in peers")
								return
							}

							sValue := fmt.Sprintf("%s", val)
							log.Debug("Peers:%s",sValue)
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
								log.Debug("Unable to get document object in settings")
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
									log.Debug("Unable to get output text area in settings keys")
									return
								}
								// sValue := fmt.Sprintf("%s:%s", name, value)

								jsonOutputTextArea.Set("innerHTML", value)
							}
						}

						default:{
							log.Debug("Default Hit")
						}
					}
				} else {
					log.Debug("Not Handled Yet")
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
