package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var client *websocket.Conn

type SMS struct {
	Destination string `json:"destination"`
	Body        string `json:"body"`
}

func webHandler(ws *websocket.Conn) {
	defer ws.Close()
	var err error

	var receivedMessage string

	if err = websocket.Message.Receive(ws, &receivedMessage); err != nil {
		fmt.Println("[ERROR] Can't receive")
		return
	}

	fmt.Println("[VERBOSE] received:" + receivedMessage)

	for i := 0; i < 2; i++ {

		sms := SMS{Destination: "6285777779927", Body: "test"}
		jsonString, _ := json.Marshal(sms)

		fmt.Println("[INFO] Send: " + string(jsonString))

		if err = websocket.Message.Send(ws, string(jsonString)); err != nil {
			fmt.Println("[ERROR] Can't send")
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(webHandler)}
		s.ServeHTTP(w, req)
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
