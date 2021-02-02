// GOOS=js GOARCH=wasm go build -o  ../assets/hive.wasm
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	// "time"
	"strconv"
)

var DNSState bool

func GetSettings() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "settings", "-g", "-j"}, splicer),
			}
			log.Debug("Settings Hit")
			val := GetData(payload, "GetSettings")
			var settings Settings
			err := json.Unmarshal(val, &settings)
			if err != nil {
				log.Error("Error in unmarshalling val in GetSettings: ", err.Error())
				return
			}
			log.Debug(settings)
			SetDisplay("Name", "innerHTML", settings.Name)

			sUsedSpace := fmt.Sprintf("%.2f %s", settings.UsedStorage*1024, "MB")
			SetDisplay("UsedSpace", "innerHTML", sUsedSpace)

			freeSpace := (settings.MaxStorage - settings.UsedStorage)
			sFreeSpace := fmt.Sprintf("%.2f %s", freeSpace*1024, "MB")
			SetDisplay("FreeSpace", "innerHTML", sFreeSpace)
			DNSState = settings.IsDNSEligible
			log.Debugf("MaxStorage: %f", settings.MaxStorage)
			SetDisplay("rangeSlider", "value", fmt.Sprintf("%s", settings.MaxStorage))
		}()
		return nil
	})
}

func GetStatus() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "status", "-j"}, splicer),
			}
			log.Debug("GetStatus Hit")
			val := GetData(payload, "GetStatus")

			var status Status
			err := json.Unmarshal(val, &status)
			if err != nil {
				log.Error("Error in unmarshalling val in GetStatus: ", err.Error())
				return
			}
			var sValue string
			if status.LoggedIn == true {
				sValue = "LoggedIn"
			} else if status.LoggedIn == false {
				sValue = "LoggedOut"
			}
			SetDisplay("LoggedIn", "innerHTML", sValue)

			// timeStamp := time.Unix(status.TotalUptimePercentage.Timestamp, 0)
			// sTimeStamp := fmt.Sprintf("%s", timeStamp.Format(time.Kitchen))
			// SetDisplay("LastConnected", "innerHTML", sTimeStamp)

		}()
		return nil
	})
}
func GetConfig() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "config", "show", "-j"}, splicer),
			}
			log.Debug("GetConfig Hit")
			val := GetData(payload, "GetConfig")
			var config Config
			err := json.Unmarshal(val, &config)
			if err != nil {
				log.Error("Error in unmarshalling val in GetConfig: ", err.Error())
				return
			}
			log.Debug(config)
			SetDisplay("SwrmPortNumber", "placeholder", config.SwarmPort)
			Attributes := make(map[string]string)
			if DNSState == false {
				Attributes["style"] = "display: none;"
				Attributes["aria-hidden"] = "true"
				Attributes["visibility"] = "hidden"
				SetMultipleDisplay("Group_62_ID", Attributes)
				return
			}
			SetDisplay("WebSocketPortNumber", "placeholder", config.WebsocketPort)
		}()
		return nil
	})
}
func SetStorageSize() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "config", "get-storage-location", "-j"}, splicer),
			}
			val, err := ModifyConfig(payload, "SetStorageSize")
			if err != nil {
				log.Error("Error in Getting Storage Location ", err.Error())
			}
			var out Out
			err = json.Unmarshal([]byte(val), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetStorageLocation: ", err.Error())
				return
			}
			log.Debug(out.Data)
		}()
		return nil
	})
}
func CheckPort(port string) (status bool, condition string) {
	if port == "" {
		return false, fmt.Sprintf("Enter A Valid Port Number")
	}
	val, err := strconv.Atoi(port)
    if err != nil {
		return false, fmt.Sprintf("Port %s is Not a Number", port)
    }
	if val < 1025 || val > 49150 {
	 	return false, fmt.Sprintf("Port %s is Unavailable", port)
	}
	return true, ""
}
func SetSwrmPortNumber() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating SwarmPort Number")
			SetDisplay("SwrmPortStatus", "innerHTML", "")
			port := GetValue("SwrmPortNumber", "value")
			Attributes := make(map[string]string)
			status, condition := CheckPort(port)
			if status == true {
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "config", "modify", "SwarmPort", port}, splicer),
			}

			log.Debugf("Payload in SetSwrmPortNumber: %s", payload)
			val, _ := ModifyConfig(payload, "SetSwrmPortNumber")
			if strings.Contains(val, "not") {
				Attributes["innerHTML"] = fmt.Sprintf("Port %s is Unavailable", port)
				Attributes["style"] = "color: red;"
				SetMultipleDisplay("SwrmPortStatus", Attributes)
				return
			}
				log.Debug("SwrmPort Updated Successfully")
				SetDisplay("SwrmPortNumber", "placeholder", port)
				SetDisplay("RestartBanner", "style", "display: block;")
				Attributes["innerHTML"] = fmt.Sprintf("SwrmPort Changed to %s", port)
				Attributes["style"] = "color: #32CD32;"
				SetMultipleDisplay("SwrmPortStatus", Attributes)
				return
			} else if status == false {
				Attributes["innerHTML"] = condition
				Attributes["style"] = "color: red;"
				SetMultipleDisplay("SwrmPortStatus", Attributes)
			}
		}()
		return nil
	})
}
func SetWebsocketPortNumber() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating SetWebsocketPortNumber Number")
			SetDisplay("WebsocketPortStatus", "innerHTML", "")
			port := GetValue("WebSocketPortNumber", "value")
			Attributes := make(map[string]string)
			status, condition := CheckPort(port)
			if status == true {
				payload := map[string]interface{}{
					"val": strings.Join([]string{"hive-cli.exe", "config", "modify", "WebsocketPort", port}, splicer),
				}
				log.Debugf("Payload in SetWebsocketPortNumber: %s", payload)
				val, _ := ModifyConfig(payload, "SetWebsocketPortNumber")
				if strings.Contains(val, "not") {
					Attributes["innerHTML"] = fmt.Sprintf("Port %s is Unavailable", port)
					Attributes["style"] = "color: red;"
					SetMultipleDisplay("WebsocketPortStatus", Attributes)
					return
				}
					log.Debug("SwrmPort Updated Successfully")
					SetDisplay("WebSocketPortNumber", "placeholder", port)
					SetDisplay("RestartBanner", "style", "display: block;")
					Attributes["innerHTML"] = fmt.Sprintf("WebsocketPort Changed to %s", port)
					Attributes["style"] = "color: #32CD32;"
					SetMultipleDisplay("WebsocketPortStatus", Attributes)
					return
		} else if status == false {
			Attributes["innerHTML"] = condition
			Attributes["style"] = "color: red;"
			SetMultipleDisplay("WebsocketPortStatus", Attributes)
		}
		}()
		return nil
	})
}
func VerifyPort() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Verifying Port Forwarding....")
			Attributes := make(map[string]string)

			Attributes["innerHTML"] = "Verifying...."
			Attributes["style"] = "color: rgba(219,219,219,1);"
			SetMultipleDisplay("PortForward", Attributes)

			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "verify-port-forward"}, splicer),
			}
			val, err := ModifyConfig(payload, "VerifyPort")
			if err != nil {
				log.Error("Error in Checking Port Forwarding Status")
				SetDisplay("PortForward", "innerHTML", "Error in Checking")
				return
			}
			log.Debugf("This is val: %s", val)

			if strings.Contains(val, "NOT") {
				log.Debug("Port Forward Not Verified")
				Attributes["innerHTML"] = "Not Forwarded &#10008;"
				Attributes["style"] = "color: rgba(244,105,50,1);"
				SetMultipleDisplay("PortForward", Attributes)
				return
			}
			log.Debug("Port Forward Verified")
			Attributes["innerHTML"] = "Port Forwarded &#10004;"
			Attributes["style"] = "color: green;"
			SetMultipleDisplay("PortForward", Attributes)
		}()
		return nil
	})
}
func ModifyStorageSize() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Changing Storage Size....")
			log.Debug(GetValue("rangeSlider", "value"))
		}()
		return nil
	})
}
