package main

// import (
//     "bufio"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	logger "github.com/ipfs/go-log/v2"
// )
//
// var log = logger.Logger("events")
//
// const (
// 	EVENTS  = "http://localhost:4343/v3/events"
// 	GATEWAY = "http://localhost:4343/v3/execute"
// )

// func UpdateStorageLocation() js.Func {
// 	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		go func() {
//
//             jsDoc := js.Global().Get("document")
// 			if !jsDoc.Truthy() {
// 				log.Error("Unable to get document object in UpdateStorageLocation")
// 				return
// 			}
// 			OutputArea := jsDoc.Call("getElementById", "SelectStoragePath")
// 			if !OutputArea.Truthy() {
// 				log.Error("Unable to get output area in UpdateStorageLocation")
// 				return
// 			}
//
//             value := fmt.Sprintf("%s", OutputArea.Get("value"))
//             log.Debug(value)
// 			// payload := map[string]interface{}{
// 			// 	"val": "hive-cli.exe%$#config%$#modify%$# %$#-j",
// 			// }
//             //
// 			// buf, err := json.Marshal(payload)
// 			// if err != nil {
// 			// 	log.Error("Error in Marshalling Payload in UpdateStorageLocation: ", err.Error())
// 			// 	return
// 			// }
// 			// resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
// 			// if err != nil {
// 			// 	log.Error("Error in getting response in UpdateStorageLocation: ", err.Error())
// 			// 	return
// 			// }
// 			// defer resp.Body.Close()
// 			// respBuf, err := ioutil.ReadAll(resp.Body)
// 			// if err != nil {
// 			// 	log.Error("Error in Reading Body in UpdateStorageLocation: ", err.Error())
// 			// 	return
// 			// }
// 			// data := make(map[string]string)
// 			// err = json.Unmarshal(respBuf, &data)
// 			// if err != nil {
// 			// 	log.Error("Error in Unmarshalling respBuf in UpdateStorageLocation: ", err.Error())
// 			// 	return
// 			// }
//             // log.Debug("This is data: ", data)
// 			// var out Out
// 			// err = json.Unmarshal([]byte(data["val"]), &out)
// 			// if err != nil {
// 			// 	log.Error("Error in Unmarshalling data in UpdateStorageLocation: ", err.Error())
// 			// 	return
// 			// }
// 			// jsDoc := js.Global().Get("document")
// 			// if !jsDoc.Truthy() {
// 			// 	log.Error("Unable to get document object in StoragePath")
// 			// 	return
// 			// }
// 			// OutputArea := jsDoc.Call("getElementById", "StoragePath")
// 			// if !OutputArea.Truthy() {
// 			// 	log.Error("Unable to get output area in StoragePath")
// 			// 	return
// 			// }
// 			// OutputArea.Set("innerHTML", out.Data)
// 		}()
// 		return nil
// 	})
// 	return jsonFunc
// }
