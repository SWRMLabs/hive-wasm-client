package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"
	// "net"
)

func GetSettings() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#settings%$#-g%$#-j",
			}
            log.Debug("Settings Hit")
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in GetSettings: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetSettings: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in Reading Body in GetSettings: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetSettings: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetSettings: ", err.Error())
				return
			}
            val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in Marshalling in GetSettings: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in GetSettings")
				return
			}
            var settings Settings
			err = json.Unmarshal(val, &settings)
			if err != nil {
				log.Error("Error in unmarshalling val in GetSettings: ", err.Error())
				return
			}
			jsDoc = js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in settings")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "Name")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in Name")
				return
			}
			OutputArea.Set("innerHTML", settings.Name)

			log.Debug(settings)
			usedSpace := (settings.UsedStorage * settings.MaxStorage) / 100
			jsDoc = js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in settings")
				return
			}
			OutputArea = jsDoc.Call("getElementById", "UsedSpace")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in UsedSpace")
				return
			}
			sUsedSpace := fmt.Sprintf("%.2f %s", usedSpace * 1024, "MB")
			OutputArea.Set("innerHTML", sUsedSpace)

			OutputArea = jsDoc.Call("getElementById", "FreeSpace")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in UsedSpace")
				return
			}
			freeSpace := (settings.MaxStorage - usedSpace)
			sFreeSpace := fmt.Sprintf("%.2f %s", freeSpace * 1024, "MB")
			OutputArea.Set("innerHTML", sFreeSpace)
		}()
		return nil
	})
	return jsonFunc
}

func GetStatus() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#status%$#-j",
			}
            log.Debug("GetStatus Hit")
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in GetStatus: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetStatus: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in Reading Body in GetStatus: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetStatus: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetStatus: ", err.Error())
				return
			}
            val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in Marshalling in GetStatus: ", err.Error())
				return
			}
			var status Status
			err = json.Unmarshal(val, &status)
			if err != nil {
				log.Error("Error in unmarshalling val in GetStatus: ", err.Error())
				return
			}
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in GetStatus")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "LoggedIn")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in LoggedIn")
				return
			}
			var sValue string
			if status.LoggedIn == true {
					sValue = "LoggedIn &#x1f7e2;"
				} else if status.LoggedIn == false {
					sValue = "LoggedOut &#10060;"
				}
			OutputArea.Set("innerHTML", sValue)

			OutputArea = jsDoc.Call("getElementById", "LastConnected")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output text area in Last Connected")
				return
			}
			timeStamp := time.Unix(status.TotalUptimePercentage.Timestamp, 0)

			sTimeStamp := fmt.Sprintf("%s", timeStamp.Format(time.Kitchen))
			OutputArea.Set("innerHTML", sTimeStamp)
		}()
		return nil
	})
	return jsonFunc
}

func GetConfig() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			payload := map[string]interface{}{
				"val": "hive-cli.exe%$#config%$#show%$#-j",
			}
            log.Debug("GetConfig Hit")
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in GetConfig: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in GetConfig: ", err.Error())
				return
			}
			defer resp.Body.Close()
			respBuf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error("Error in Reading Body in GetConfig: ", err.Error())
				return
			}
			data := make(map[string]string)
			err = json.Unmarshal(respBuf, &data)
			if err != nil {
				log.Error("Error in Unmarshalling respBuf in GetConfig: ", err.Error())
				return
			}

			var out Out
			err = json.Unmarshal([]byte(data["val"]), &out)
			if err != nil {
				log.Error("Error in Unmarshalling data in GetConfig: ", err.Error())
				return
			}
            val, err := json.MarshalIndent(out.Data, "", " ")
			if err != nil {
				log.Error("Error in Marshalling in GetConfig: ", err.Error())
				return
			}
			log.Debug(data["val"])
			var config Config
			err = json.Unmarshal(val, &config)
			if err != nil {
				log.Error("Error in unmarshalling val in GetConfig: ", err.Error())
				return
			}
			log.Debug(config)
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in GetConfig")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "SwrmPortNumber")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in SwrmPortNumber")
				return
			}
			OutputArea.Set("placeholder", config.SwarmPort)

			OutputArea = jsDoc.Call("getElementById", "WebSocketPortNumber")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in WebSocketPortNumber")
				return
			}
			OutputArea.Set("placeholder", config.WebsocketPort)
		}()
		return nil
	})
	return jsonFunc
}

func SetSwrmPortNumber() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating SwarmPort Number")
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in SetSwrmPortNumber")
				return
			}
			OutputArea := jsDoc.Call("getElementById", "SwrmPortStatus")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in SwrmPortStatus")
				return
			}
			OutputArea.Set("innerHTML", "")

			OutputArea = jsDoc.Call("getElementById", "SwrmPortNumber")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in SetSwrmPortNumber")
				return
			}

			port := fmt.Sprintf("%s", OutputArea.Get("value"))
			// host := "127.0.0.1"
			// log.Debugf("Checking Availability of Port: %s", port)
			//
			// ln, err := net.Listen("tcp", "127.0.0.1:" + port)
			// if err != nil {
			// 	log.Debug("server")
			// 	log.Error(err.Error())
			// }
			// defer ln.Close()
			//
			//
			// log.Debugf("Listening on %s", net.JoinHostPort(host, port))
			// conn, err := net.DialTCP("tcp", nil, net.JoinHostPort(host, port))
			// if err != nil {
			// 	log.Debug("client")
			// 	log.Error(err.Error())
			// 	OutputArea := jsDoc.Call("getElementById", "SwrmPortStatus")
			// 	if !OutputArea.Truthy() {
			// 		log.Error("Unable to get output area in SwrmPortStatus")
			// 		return
			// 	}
			// 	OutputArea.Set("innerHTML", "This Port is unavailable. Please Select a different Port.")
			// 	OutputArea.Set("style", "color:rgba(255,15,15,1)")
			// 	return
			// 	}
			// 	conn.Close()
			// 	ln.Close()
			// log.Debugf("Port:%s is available", port)

			sValue := "hive-cli.exe%$#config%$#modify%$#SwarmPort%$#" + port
			log.Debug(sValue)
			payload := map[string]interface{}{
				"val": sValue,
			}
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in SetSwrmPortNumber: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in SetSwrmPortNumber: ", err.Error())
				return
			}
			defer resp.Body.Close()
			log.Debug("This is response from SetSwrmPortNumber: ",resp)

			payload = map[string]interface{}{
				"val": "hive-cli.exe%$#settings%$#-j",
			}

			buf, err = json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in SetSwrmPortNumber: ", err.Error())
				return
			}
			resp, err = http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in SetSwrmPortNumber: ", err.Error())
				return
			}
			defer resp.Body.Close()
			log.Debug("This is response from SetSwrmPortNumber: ",resp)

			if err == nil {
				log.Debug("SwrmPort Updated Successfully")
				OutputArea.Set("placeholder", port)
			}

		}()
		return nil
	})
	return jsonFunc
}

func SetWebsocketPortNumber() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			log.Debug("Updating WebSocketPort Number")
			jsDoc := js.Global().Get("document")
			if !jsDoc.Truthy() {
				log.Error("Unable to get document object in WebsocketPortNumber")
				return
			}

			OutputArea := jsDoc.Call("getElementById", "WebSocketPortNumber")
			if !OutputArea.Truthy() {
				log.Error("Unable to get output area in WebSocketPortNumber")
				return
			}
			value := OutputArea.Get("value")
			log.Debug(value)
			sep := "%$#"
			sValue := "hive-cli.exe%$#config%$#modify%$#WebsocketPort" + sep + fmt.Sprintf("%s",value)
			log.Debug(sValue)
			payload := map[string]interface{}{
				"val": sValue,
			}
			buf, err := json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in WebSocketPortNumber: ", err.Error())
				return
			}
			resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in WebSocketPortNumber: ", err.Error())
				return
			}
			defer resp.Body.Close()
			log.Debug("This is response from WebSocketPortNumber: ",resp)

			payload = map[string]interface{}{
				"val": "hive-cli.exe%$#settings%$#-j",
			}

			buf, err = json.Marshal(payload)
			if err != nil {
				log.Error("Error in Marshalling Payload in WebSocketPortNumber: ", err.Error())
				return
			}
			resp, err = http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Error("Error in getting response in WebSocketPortNumber: ", err.Error())
				return
			}
			defer resp.Body.Close()
			log.Debug("This is response from WebSocketPortNumber: ",resp)

			if err == nil {
				log.Debug("WebSocketPort Updated Successfully")
				OutputArea.Set("placeholder", value)
			}
		}()
		return nil
	})
	return jsonFunc
}
