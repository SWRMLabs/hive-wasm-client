// GOOS=js GOARCH=wasm go build -o  ../assets/hive.wasm

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	logger "github.com/ipfs/go-log/v2"
	"net/http"
	"reflect"
	"strings"
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
						log.Debug("Error encountered in Marshalling")
						return
					}

					if (event.Result.Topic == "Status") || (event.Result.Topic == "Balance") || (event.Result.Topic == "BalanceCycle") || (event.Result.Topic == "Peers") || (event.Result.Topic == "Settings") {
						switch event.Result.Topic {

						case "Status":
							{
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

								for _, task := range status.TaskManagerStatus {
									sName := task.Name
									if sName == "Idle" {
										continue
									}
									OutputArea := jsDoc.Call("createElement", "div")
									if !OutputArea.Truthy() {
										log.Debug("Unable to get output area task manager name")
										return
									}
									OutputArea.Set("innerHTML", sName)
									jsDoc.Call("getElementById", "taskmanagerstatusname").Call("appendChild", OutputArea)

									sStatus := task.Status
									OutputArea = jsDoc.Call("createElement", "div")
									if !OutputArea.Truthy() {
										log.Debug("Unable to get output area task manager status")
										return
									}
									OutputArea.Set("innerHTML", sStatus)
									jsDoc.Call("getElementById", "taskmanagerstatusstatus").Call("appendChild", OutputArea)

									sAdditionalStatus := task.AdditionalStatus
									OutputArea = jsDoc.Call("createElement", "div")
									if !OutputArea.Truthy() {
										log.Debug("Unable to get output area task manager additional status")
										return
									}
									OutputArea.Set("innerHTML", sAdditionalStatus)
									jsDoc.Call("getElementById", "taskmanagerstatusAS").Call("appendChild", OutputArea)
								}

								serverStatus := reflect.ValueOf(&status.ServerDetails).Elem()

								for key := 0; key < serverStatus.NumField(); key++ {
									name := serverStatus.Type().Field(key).Name
									value := serverStatus.Field(key).Interface()
									log.Debug(name)
									log.Debug(value)
									OutputArea := jsDoc.Call("getElementById", name)
									if !OutputArea.Truthy() {
										log.Debug("Unable to get output text area in server status")
										return
									}
									OutputArea.Set("innerHTML", value)
								}

								values := reflect.ValueOf(&status).Elem()

								for key := 0; key < values.NumField(); key++ {
									name := values.Type().Field(key).Name
									value := values.Field(key).Interface()

									if (name == "TaskManagerStatus") || (name == "TotalUptimePercentage") || (name == "SessionStartTime") || (name == "ServerDetails") {
										continue
									}

									log.Debug(values.Type().Field(key).Name)
									OutputArea := jsDoc.Call("getElementById", name)
									if !OutputArea.Truthy() {
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
									} else if value == false {
										switch name {
										case "LoggedIn":
											sValue = "LoggedOut"

										case "DaemonRunning":
											sValue = "OFFLINE"
										}
									}
									OutputArea.Set("innerHTML", sValue)

								}

								OutputArea := jsDoc.Call("getElementById", "Time")
								if !OutputArea.Truthy() {
									log.Debug("Unable to get output text area in Time")
									return
								}
								time := status.TotalUptimePercentage.SecondsFromInception

								days := (time) / (24 * 3600)
								time = time % (24 * 3600)
								hours := time / 3600
								time = time % 3600
								minutes := time / 60
								time = time % 60
								seconds := time
								sValue := fmt.Sprintf("%d:%d:%d:%d", days, hours, minutes, seconds)
								log.Debug(sValue)
								OutputArea.Set("innerHTML", sValue)

								OutputArea = jsDoc.Call("getElementById", "percentageNumber")
								if !OutputArea.Truthy() {
									log.Debug("Unable to get output text area in percentagenumber")
									return
								}
								sFloat := fmt.Sprintf("%.2f", status.TotalUptimePercentage.Percentage)
								sValue = fmt.Sprintf("%s %s", sFloat, "%")
								OutputArea.Set("innerHTML", sValue)

								// OutputArea = jsDoc.Call("getElementById", "secondsfrominception")
								// if !OutputArea.Truthy() {
								// 	log.Debug("Unable to get output text area in secondsfromincepti")
								// 	return
								// }
								// OutputArea.Set("innerHTML", status.TotalUptimePercentage.SecondsFromInception)
							}

						case "Balance":
							{
								log.Debug("Balance Hit")

								jsDoc := js.Global().Get("document")
								if !jsDoc.Truthy() {
									log.Debug("Unable to get document object in balance")
									return
								}

								OutputArea := jsDoc.Call("getElementById", "confirmedBalance")
								if !OutputArea.Truthy() {
									log.Debug("Unable to get output text area in balance")
									return
								}
								sFloat := fmt.Sprintf("%s", val)
								for i, value := range sFloat {
									if strings.ContainsAny(string(value), ".") && (i+3) <= len(sFloat) {
										sFloat = sFloat[0:i+1] + sFloat[i+1:i+3]
										break
									}
								}
								sValue := fmt.Sprintf("%s %s", sFloat, "SWRM")
								OutputArea.Set("innerHTML", sValue)
							}

						case "BalanceCycle":
							{
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
									log.Debug("Unable to get document object in balance cycle")
									return
								}

								OutputArea := jsDoc.Call("getElementById", "Pending")
								if !OutputArea.Truthy() {
									log.Debug("Unable to get output text area in balance cycle keys")
									return
								}

								sValue := fmt.Sprintf("%f %s", (bcnBalance.Owned - bcnBalance.Owe), "SWRM")
								log.Debug("This is Pending:", sValue)
								OutputArea.Set("innerHTML", sValue)

							}

						case "Peers":
							{
								log.Debug("Peers Hit")

								jsDoc := js.Global().Get("document")
								if !jsDoc.Truthy() {
									log.Debug("Unable to get document object in peers")
									return
								}

								OutputArea := jsDoc.Call("getElementById", "PeersData")
								if !OutputArea.Truthy() {
									log.Debug("Unable to get output text area in peers")
									return
								}

								sValue := fmt.Sprintf("%s", val)
								log.Debug("Peers:%s", sValue)
								OutputArea.Set("innerHTML", sValue)

							}

						case "Settings":
							{
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
									if (name == "PeerID") || (name == "Role") || (name == "MaxStorage") || (name == "UsedStorage") {
										log.Debug(value)
										OutputArea := jsDoc.Call("getElementById", name)
										if !OutputArea.Truthy() {
											log.Debug("Unable to get output text area in settings keys")
											return
										}
										sValue := value
										if name == "MaxStorage" {
											sValue = fmt.Sprintf("%.2f %s", value, "GB")
										}
										if name == "UsedStorage" {
											sValue = fmt.Sprintf("%.2f %s", value, "%")
										}

										OutputArea.Set("innerHTML", sValue)
									}

								}
							}
						default:
							{
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
