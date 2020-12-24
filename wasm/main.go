// GOOS=js GOARCH=wasm go build -o  ../assets/hive.wasm

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/ipfs/go-log/v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"syscall/js"
	"time"
)

var log = logger.Logger("events")

const (
	EVENTS  = "http://localhost:4343/v3/events"
	GATEWAY = "http://localhost:4343/v3/execute"
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
				var eventsDataString string

				line, isPrefix, err := reader.ReadLine()
				if err != nil {
					log.Debugf("Error in reading data string: %s", err.Error())
					continue
				}
				eventsDataString = string(line)

				for isPrefix {
					line, isPrefix, err = reader.ReadLine()
					if err != nil {
						log.Debugf("Error in reading prefixed data string: %s", err.Error())
						continue
					}
					eventsDataString += string(line)
				}
				log.Debug(eventsDataString)
				var event Event
				err = json.Unmarshal([]byte(eventsDataString), &event)
				if err != nil {
					log.Error("Error in Unmarshalling eventsDataString:", err.Error())
					return
				}
				var out Out
				err = json.Unmarshal([]byte(event.Result.Val), &out)
				if err != nil {
					log.Error("Error in Unmarshalling Out in %s : ", event.Result.Val, err.Error())
					return
				}
				val, err := json.MarshalIndent(out.Data, "", " ")
				if err != nil {
					log.Error("Error encountered in Marshalling: ", err.Error())
					return
				}

				if (event.Result.Topic == "Status") || (event.Result.Topic == "Balance") || (event.Result.Topic == "BalanceCycle") || (event.Result.Topic == "Peers") || (event.Result.Topic == "Settings") || (event.Result.Topic == "Settlement") || (event.Result.Topic == "Earning") {
					switch event.Result.Topic {
					case "Status":
						{
							log.Debug("Status Hit")
							var status Status
							err = json.Unmarshal(val, &status)
							if err != nil {
								log.Error("Error in Unmarshalling Status:", err.Error())
								return
							}
							log.Debug("This is Status: ", status)
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in status")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "taskmanagerstatusname")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output area task manager name")
								return
							}
							OutputArea.Set("innerHTML", "")

							OutputArea = jsDoc.Call("getElementById", "taskmanagerstatusstatus")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output area task manager status")
								return
							}
							OutputArea.Set("innerHTML", "")

							OutputArea = jsDoc.Call("getElementById", "taskmanagerstatusAS")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output area task manager Additional Status")
								return
							}
							OutputArea.Set("innerHTML", "")

							for _, task := range status.TaskManagerStatus {
								sName := task.Name
								if sName == "Idle" {
									continue
								}
								OutputArea := jsDoc.Call("createElement", "div")
								if !OutputArea.Truthy() {
									log.Error("Unable to get create div in task manager name")
									return
								}
								OutputArea.Set("innerHTML", sName)
								jsDoc.Call("getElementById", "taskmanagerstatusname").Call("appendChild", OutputArea)

								sStatus := task.Status
								OutputArea = jsDoc.Call("createElement", "div")
								if !OutputArea.Truthy() {
									log.Error("Unable to get create div in task manager status")
									return
								}
								OutputArea.Set("innerHTML", sStatus)
								jsDoc.Call("getElementById", "taskmanagerstatusstatus").Call("appendChild", OutputArea)

								sAdditionalStatus := task.AdditionalStatus
								OutputArea = jsDoc.Call("createElement", "div")
								if !OutputArea.Truthy() {
									log.Error("Unable to get create div in task manager additional status")
									return
								}
								OutputArea.Set("innerHTML", sAdditionalStatus)
								jsDoc.Call("getElementById", "taskmanagerstatusAS").Call("appendChild", OutputArea)
							}

							serverStatus := reflect.ValueOf(&status.ServerDetails).Elem()

							for key := 0; key < serverStatus.NumField(); key++ {
								name := serverStatus.Type().Field(key).Name
								value := serverStatus.Field(key).Interface()
								OutputArea := jsDoc.Call("getElementById", name)
								if !OutputArea.Truthy() {
									log.Error("Unable to get output text area in server status")
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
								OutputArea := jsDoc.Call("getElementById", name)
								if !OutputArea.Truthy() {
									log.Error("Unable to get output text area in status keys")
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
							OutputArea = jsDoc.Call("getElementById", "LastConnected")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in Last Connected")
								return
							}
							timeStamp := time.Unix(status.TotalUptimePercentage.Timestamp, 0)

							sTimeStamp := fmt.Sprintf("%s", timeStamp.Format(time.Kitchen))
							OutputArea.Set("innerHTML", sTimeStamp)

							OutputArea = jsDoc.Call("getElementById", "Time")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in Time")
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
							sValue := fmt.Sprintf("%ddays %dhours %dmins %dsecs", days, hours, minutes, seconds)
							OutputArea.Set("innerHTML", sValue)

							OutputArea = jsDoc.Call("getElementById", "percentageNumber")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in percentagenumber")
								return
							}
							sFloat := fmt.Sprintf("%.2f", status.TotalUptimePercentage.Percentage)
							sValue = fmt.Sprintf("%s %s", sFloat, "%")
							OutputArea.Set("innerHTML", sValue)
						}

					case "Balance":
						{
							log.Debug("Balance Hit")

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in balance")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "confirmedBalance")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in balance")
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
							log.Debugf("This is Main Balance: %s", sValue)
							OutputArea.Set("innerHTML", sValue)
						}

					case "Settlement":
						{
							log.Debug("Settlement Hit")
							var settlement Settlement
							err = json.Unmarshal(val, &settlement)
							if err != nil {
								log.Error("Error Unmarshalling settlement: ", err.Error())
								return
							}
							log.Debug("This is Settlement: ", settlement)
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in settlement")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "NextDistribution")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in settlement")
								return
							}
							date := (settlement.Date).Format("02-01-2006")
							time := (settlement.Date).Format(time.Kitchen)
							sDateTime := fmt.Sprintf("%s %s", date, time)
							OutputArea.Set("innerHTML", sDateTime)
						}
					case "BalanceCycle":
						{
							log.Debug("BCN Hit")
							var bcnBalance BCNBalance
							err = json.Unmarshal(val, &bcnBalance)
							if err != nil {
								log.Error("Error in Unmarshalling BCN Balance:", err.Error())
								return
							}
							log.Debug("This is Balance Cycle: ", bcnBalance)
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in balance cycle")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "Pending")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in Pending")
								return
							}
							sValue := fmt.Sprintf("%f %s", (bcnBalance.Owned - bcnBalance.Owe), "SWRM")
							OutputArea.Set("innerHTML", sValue)
						}

					case "Peers":
						{
							log.Debug("Peers Hit")
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in peers")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "PeersData")
							if !OutputArea.Truthy() {
								log.Error("Unable to get output text area in peers")
								return
							}
							sValue := fmt.Sprintf("%s", val)
							log.Debugf("This is Number of Peers: %s", val)
							OutputArea.Set("innerHTML", sValue)

						}

					case "Earning":
						{
							log.Debug("Earning Hit")
							var netEarnings NetEarnings
							err = json.Unmarshal(val, &netEarnings)
							if err != nil {
								log.Error("Error in Unmarshalling Net Earnings: ", err.Error())
								return
							}
							log.Debugf("%+v",netEarnings)

							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in Earning")
								return
							}
							OutputArea := jsDoc.Call("getElementById", "DevicesDropDown")
							if !OutputArea.Truthy() {
								log.Error("Unable to get create div in task manager additional status")
								return
							}
							OutputArea.Set("innerHTML", "")
							OutputArea = jsDoc.Call("createElement", "option")
							if !OutputArea.Truthy() {
								log.Error("Unable to get create div in task manager additional status")
								return
							}
							OutputArea.Set("innerHTML", "ALL DEVICES")
							OutputArea.Set("value", "ALL DEVICES")
							jsDoc.Call("getElementById", "DevicesDropDown").Call("appendChild", OutputArea)

							for _, value := range(netEarnings.Devices) {

								OutputArea := jsDoc.Call("createElement", "option")
								if !OutputArea.Truthy() {
									log.Error("Unable to get create div in task manager additional status")
									return
								}
								sOption := fmt.Sprintf("%s-%s",value.PeerId, value.Name)
								OutputArea.Set("innerHTML", sOption)
								OutputArea.Set("value", value.PeerId)
								jsDoc.Call("getElementById", "DevicesDropDown").Call("appendChild", OutputArea)
							}

							log.Debugf("This is Device Total: %+v ", netEarnings.DeviceTotal)

						}


					case "Settings":
						{
							log.Debug("Settings Hit")
							var settings Settings
							err = json.Unmarshal(val, &settings)
							if err != nil {
								log.Error("Error in Unmarshalling Settings: ", err.Error())
								return
							}
							log.Debug("This is Settings: ", settings)
							jsDoc := js.Global().Get("document")
							if !jsDoc.Truthy() {
								log.Error("Unable to get document object in settings")
								return
							}

							values := reflect.ValueOf(&settings).Elem()

							for key := 0; key < values.NumField(); key++ {
								name := values.Type().Field(key).Name
								value := values.Field(key).Interface()
								if (name == "MaxStorage") || (name == "UsedStorage") {
									OutputArea := jsDoc.Call("getElementById", name)
									if !OutputArea.Truthy() {
										log.Error("Unable to get output text area in settings keys")
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

			}
		}()
		return nil
	})
	return jsonfunc
}

func GetID() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#id%$#-j",
			}
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in marshalling payload in GetID: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetID: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading body in GetID: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetID: ", err.Error())
				return
			}
			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetID", err.Error())
				return
			}
			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in Marshalling out in GetID: ", err.Error())
				return
			}
			var id ID
			err = json.Unmarshal(val, &id)
			if err != nil {
				log.Error("Error in Unmarshalling ID in GetID: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in ID")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Address")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in Address")
				return
			}
			OutputArea.Set("innerHTML", "")

			for _, value := range id.Addresses {
				OutputArea := jsDoc.Call("createElement", "div")
				if !OutputArea.Truthy() {
					log.Error("Unable to create div in Address")
					return
				}
				OutputArea.Set("innerHTML", value)
				jsDoc.Call("getElementById", "Address").Call("appendChild", OutputArea)
				OutputArea = jsDoc.Call("createElement", "br")
				if !OutputArea.Truthy() {
					log.Debug("Unable to create break in Address")
					return
				}
				jsDoc.Call("getElementById", "Address").Call("appendChild", OutputArea)
			}

			OutputArea = jsDoc.Call("getElementById", "PeerID")
			if !OutputArea.Truthy() {
				log.Debug("Unable to get output area in PeerID")
				return
			}
			OutputArea.Set("innerHTML", id.PeerID)
		}()
		return nil
	})
	return jsonFunc
}

func GetPeers() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#swarm%$#peers%$#-j",
			}
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error(err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetPeers: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading body in GetPeers: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetPeers: ", err.Error())
				return
			}
			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetPeers: ", err)
				return
			}
			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in Marshalling in GetPeers: ", err.Error())
				return
			}
			var swarmPeers []string
			err = json.Unmarshal(val, &swarmPeers)
			if err != nil {
				log.Error("Error in Unmarshalling SwarmPeers: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Debug("Unable to get document object in Peers")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Peers")
			if !OutputArea.Truthy() {
				log.Debug("Unable to get output area in Peers")
				return
			}
			OutputArea.Set("innerHTML", "")

			for _, value := range swarmPeers {
				OutputArea := jsDoc.Call("createElement", "div")
				if !OutputArea.Truthy() {
					log.Debug("Unable to get create div in Peers")
					return
				}
				OutputArea.Set("innerHTML", value)
				jsDoc.Call("getElementById", "Peers").Call("appendChild", OutputArea)
				OutputArea = jsDoc.Call("createElement", "br")
				if !OutputArea.Truthy() {
					log.Debug("Unable to get create break in Peers")
					return
				}
				jsDoc.Call("getElementById", "Peers").Call("appendChild", OutputArea)
			}
		}()
		return nil
	})
	return jsonFunc
}

func GetStorageLocation() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#config%$#get-storage-location%$#-j",
			}

			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in GetStorageLocation: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetStorageLocation: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in Reading Body in GetStorageLocation: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetStorageLocation: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetStorageLocation: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in StoragePath")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "StoragePath")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in StoragePath")
				return
			}
			OutputArea.Set("innerHTML", out.Data)
		}()
		return nil
	})
	return jsonFunc
}

func GetProfile() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#profile%$#-j",
			}
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in GetProfile: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetProfile: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading respBuf in GetProfile: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in unmarshalling respBuf in GetProfile: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in unmarshalling data in GetProfile: ", err.Error())
				return
			}
			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in marshalling out in GetProfile: ", err.Error())
				return
			}
			var profile Profile
			err = json.Unmarshal(val, &profile)
			if err != nil {
				log.Error("Error in unmarshalling val in GetProfile: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in Profile")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Email")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in Email")
				return
			}
			OutputArea.Set("innerHTML", profile.Email)

			OutputArea = jsDoc.Call("getElementById", "Role")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in Role")
				return
			}
			OutputArea.Set("innerHTML", profile.Role)

		}()
		return nil
	})
	return jsonFunc
}

func GetBandwidth() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#stat%$#bandwidth%$#-j",
			}

			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in marshalling Payload in GetBandwidth: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetBandwidth: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading Body in GetBandwidth: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in unmarshalling respBuf in GetBandwidth: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in unmarshalling data in GetBandwidth: ", err.Error())
				return
			}
			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in marshalling out in GetBandwidth: ", err.Error())
				return
			}
			var bandwidth Bandwidth
			err = json.Unmarshal(val, &bandwidth)
			if err != nil {
				log.Error("Error in unmarshalling val in GetBandwidth: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in Bandwidth")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Incoming")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in incoming")
				return
			}
			sIncoming := fmt.Sprintf("%.3f %s", (bandwidth.Incoming)/1000, "MB")
			OutputArea.Set("innerHTML", sIncoming)

			OutputArea = jsDoc.Call("getElementById", "Outgoing")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in outgoing")
				return
			}
			sOutgoing := fmt.Sprintf("%.3f %s", (bandwidth.Outgoing)/1000, "MB")
			OutputArea.Set("innerHTML", sOutgoing)
		}()
		return nil
	})
	return jsonFunc
}

func GetEarning() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#earning%$#-g%$#-j",
			}

			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in marshalling Payload in GetEarning: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetEarning: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading Body in GetEarning: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in unmarshalling respBuf in GetEarning: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in unmarshalling data in GetEarning: ", err.Error())
				return
			}
			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in marshalling out in GetBandwidth: ", err.Error())
				return
			}
			var netEarnings NetEarnings
			err = json.Unmarshal(val, &netEarnings)
			if err != nil {
				log.Error("Error in Unmarshalling Net Earnings in GetEaring: ", err.Error())
				return
			}

			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in Bandwidth")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "DevicesDropDown")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in GetEarning")
				return
			}
			value := fmt.Sprintf("%s", OutputArea.Get("value"))
			if value == "ALL DEVICES" {

				earned := fmt.Sprintf("%.5f %s",netEarnings.DeviceTotal.Earned, "SWRM")
				download := fmt.Sprintf("%.0f %s",(netEarnings.DeviceTotal.Download)/1000, "MB")
				served := fmt.Sprintf("%.2f %s",(netEarnings.DeviceTotal.Served)/1000, "MB")

				OutputArea := jsDoc.Call("getElementById", "EarnedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in EarnedCycle")
					return
				}
				OutputArea.Set("innerHTML", earned)

				OutputArea = jsDoc.Call("getElementById", "DownloadedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in DownloadedCycle")
					return
				}
				OutputArea.Set("innerHTML", download)

				OutputArea = jsDoc.Call("getElementById", "ServedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in ServedCycle")
					return
				}
				OutputArea.Set("innerHTML", served)

			} else if value != "" {
				earned := fmt.Sprintf("%.5f %s",netEarnings.Data[value][0].Earned, "SWRM")
				download := fmt.Sprintf("%.2f %s",(netEarnings.Data[value][0].Download)/1000, "MB")
				served := fmt.Sprintf("%.2f %s",(netEarnings.Data[value][0].Served)/1000, "MB")

				OutputArea := jsDoc.Call("getElementById", "EarnedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in EarnedCycle")
					return
				}
				OutputArea.Set("innerHTML", earned)

				OutputArea = jsDoc.Call("getElementById", "DownloadedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in DownloadedCycle")
					return
				}
				OutputArea.Set("innerHTML", download)

				OutputArea = jsDoc.Call("getElementById", "ServedCycle")
				if !OutputArea.Truthy() {
					log.Error("Unable to get output area in ServedCycle")
					return
				}
				OutputArea.Set("innerHTML", served)
			}
		}()
		return nil
	})
	return jsonFunc
}

func GetVersion() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#version%$#-j",
			}

			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in marshalling payload in GetVersion: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetVersion: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in reading respBuf in GetVersion: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in unmarshalling respbuf in GetVersion: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in unmarshalling data in GetVersion: ", err.Error())
				return
			}

			val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in marshalling out in GetVersion: ", err.Error())
				return
			}
			var version Version
			err = json.Unmarshal(val, &version)
			if err != nil {
				log.Error("Error in unmarshalling val in GetVersion: ", err.Error())
				return
			}

			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in Version")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Version")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in Version")
				return
			}
			OutputArea.Set("innerHTML", version.AppVersion)
		}()
		return nil
	})
	return jsonFunc
}

func main() {
	logger.SetLogLevel("*", "Debug")
	js.Global().Set("GetVersion", GetVersion())
	js.Global().Set("GetProfile", GetProfile())
	js.Global().Set("GetEarning", GetEarning())
	js.Global().Set("GetBandwidth", GetBandwidth())
	js.Global().Set("GetStorageLocation", GetStorageLocation())
	js.Global().Set("GetID", GetID())
	js.Global().Set("GetPeers", GetPeers())
	js.Global().Set("Events", Events())
	<-make(chan bool)
}
