package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"syscall/js"
	"time"
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

			usedSpace := (settings.UsedStorage * settings.MaxStorage) / 100
			sUsedSpace := fmt.Sprintf("%.2f %s", usedSpace*1024, "MB")
			SetDisplay("UsedSpace", "innerHTML", sUsedSpace)

			freeSpace := (settings.MaxStorage - usedSpace)
			sFreeSpace := fmt.Sprintf("%.2f %s", freeSpace*1024, "MB")
			SetDisplay("FreeSpace", "innerHTML", sFreeSpace)
			DNSState = settings.IsDNSEligible

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

			timeStamp := time.Unix(status.TotalUptimePercentage.Timestamp, 0)
			sTimeStamp := fmt.Sprintf("%s", timeStamp.Format(time.Kitchen))
			SetDisplay("LastConnected", "innerHTML", sTimeStamp)

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
			if DNSState == true {
				Attributes["disabled"] = "false"
			} else if DNSState == false {
				Attributes["disabled"] = "true"
			}
			Attributes["placeholder"] = config.WebsocketPort
			SetMultipleDisplay("WebSocketPortNumber", Attributes)
			SetMultipleDisplay("WebSocketPortNumberButton", Attributes)

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
			// path := fmt.Sprintf("%s", out.Data)
			// usage := du.NewDiskUsage(path)
			// log.Debug("Free:", usage.Free())
			// log.Debug("Available:", usage.Available())
			// log.Debug("Size:", usage.Size())
			// log.Debug("Used:", usage.Used())
			// log.Debug("Usage:", usage.Usage()*100)
		}()
		return nil
	})
}
func Check(port string) (status bool, err error) {
	// Concatenate a colon and the port
	log.Debug("Checking port availability")
	host := ":" + port
	// Try to create a server with the port
	server, err := net.Listen("tcp", host)
	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}
	// close the server
	server.Close()
	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil
}
func SetSwrmPortNumber() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating SwarmPort Number")
			SetDisplay("SwrmPortStatus", "innerHTML", "")
			port := GetValue("SwrmPortNumber", "value")

			log.Debug("Port Available")
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "config", "modify", "SwarmPort", port}, splicer),
			}
			log.Debugf("Payload in SetSwrmPortNumber: %s", payload)
			_, err := ModifyConfig(payload, "SetSwrmPortNumber")
			if err != nil {
				log.Error("Error in Modifying Config in SetSwrmPortNumber ", err.Error())
			} else if err == nil {
				log.Debug("SwrmPort Updated Successfully")
				SetDisplay("SwrmPortNumber", "placeholder", port)
				SetDisplay("RestartBanner", "display", "block")
			}
		}()
		return nil
	})
}
func SetWebsocketPortNumber() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating SetWebsocketPortNumber Number")
			SetDisplay("SwrmPortStatus", "innerHTML", "")
			port := GetValue("WebSocketPortNumber", "value")

			log.Debug("Port Available")
			payload := map[string]interface{}{
				"val": strings.Join([]string{"hive-cli.exe", "config", "modify", "WebsocketPort", port}, splicer),
			}
			log.Debugf("Payload in SetWebsocketPortNumber: %s", payload)
			_, err := ModifyConfig(payload, "SetWebsocketPortNumber")
			if err != nil {
				log.Error("Error in Modifying Config in SetWebsocketPortNumber ", err.Error())
			} else if err == nil {
				log.Debug("WebsocketPort Updated Successfully")
				SetDisplay("WebSocketPortNumber", "placeholder", port)
				SetDisplay("RestartBanner", "display", "block")
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
				Attributes["style"] = "color: red;"
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
