package main

//go:generate GOOS=js GOARCH=wasm go build -o  ../assets/hive.wasm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"syscall/js"

	logger "github.com/ipfs/go-log/v2"
)

var log = logger.Logger("sock/server")

const (
	GATEWAY = "http://localhost:4343/v3/execute"
)

func ID() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#id%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()
				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})
	return jsonFunc
}

func Status() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#status%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Config() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#config%$#show%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Peers() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#swarm%$#peers%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Profile() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#profile%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func MainBalance() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#balance%$#-t%$#all%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func SettlementBalance() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#balance%$#-t%$#settlement%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func CycleBalance() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#balance%$#-t%$#cycle%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Earning() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#earning%$#-g%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Settings() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#settings%$#-g%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func Version() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {
				payload := map[string]interface{}{
					"val": "hive-cli.exe%$#version%$#-j",
				}

				buf, err := json.Marshal(payload)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}

				resp, err := http.Post(GATEWAY, "application/json", bytes.NewReader(buf))
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				defer resp.Body.Close()

				respBuf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
					reject.Invoke("Failed")
					return
				}
				data := make(map[string]string)
				json.Unmarshal(respBuf, &data)
				log.Debug(data["val"])
				resolve.Invoke(data["val"])
			}()

			return nil

		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
		return nil
	})

	return jsonFunc
}

func main() {
	logger.SetLogLevel("*", "Debug")

	js.Global().Set("Id", ID())
	js.Global().Set("Status", Status())
	js.Global().Set("Peers", Peers())
	js.Global().Set("Config", Config())
	js.Global().Set("Profile", Profile())
	js.Global().Set("MainBalance", MainBalance())
	js.Global().Set("SettlementBalance", SettlementBalance())
	js.Global().Set("CycleBalance", CycleBalance())
	js.Global().Set("Earning", Earning())
	js.Global().Set("Settings", Settings())
	js.Global().Set("Version", Version())
	<-make(chan bool)
}
