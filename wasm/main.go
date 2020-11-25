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
					"val": "hive%$#id%$#-j",
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
				log.Debug(string(respBuf))
				resolve.Invoke(string(respBuf))
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

	js.Global().Set("id", ID())
	<-make(chan bool)
}