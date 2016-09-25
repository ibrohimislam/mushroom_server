package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/websocket"
)

var (
	db          *sql.DB
	queryGet    *sql.Stmt
	queryUpdate *sql.Stmt
)

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

	for {

		var id string
		var nomor string
		var pesan string
		var status int

		err := queryGet.QueryRow().Scan(&id, &nomor, &pesan, &status)

		switch {
		case err == sql.ErrNoRows:
			// do nothing
		case err != nil:
			panic(err)
		default:
			sms := SMS{Destination: nomor, Body: pesan}
			jsonString, _ := json.Marshal(sms)

			fmt.Println("[INFO] Send: " + string(jsonString))

			_, err = queryUpdate.Exec(id)
			if err != nil {
				fmt.Println("[ERROR] SQL exec error: " + err.Error())
			}

			if err = websocket.Message.Send(ws, string(jsonString)); err != nil {
				fmt.Println("[ERROR] Can't send")
				break
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	var err error

	db, err = sql.Open("mysql", "root@tcp(172.17.0.2:3306)/db_risthanata")
	if err != nil {
		panic("SQL connect error: " + err.Error())
	}

	queryGet, err = db.Prepare("SELECT * FROM tb_smsg WHERE status = 0 LIMIT 1")
	if err != nil {
		panic("SQL prepare error: " + err.Error())
	}

	queryUpdate, err = db.Prepare("UPDATE tb_smsg SET status = 1 WHERE id = ?")
	if err != nil {
		panic("SQL prepare error: " + err.Error())
	}

	if err != nil {
		panic("DB Connection Error: " + err.Error())
	}

	fmt.Println("[INFO] Server started." + string(jsonString))

	http.HandleFunc("/smsgateway", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(webHandler)}
		s.ServeHTTP(w, req)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
