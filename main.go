package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type jsonMsg struct {
	UserName  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws/v1", ws1)
	http.HandleFunc("/ws/v2", ws2)
	http.HandleFunc("/ws/v3", ws3)
	http.HandleFunc("/ws/v4", ws4)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Print on the client
func ws1(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, nil)
	go func(conn *websocket.Conn) {
		mType, msg, _ := conn.ReadMessage()
		conn.WriteMessage(mType, msg)
	}(conn)
}

// Print on the server
func ws2(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, nil)
	go func(conn *websocket.Conn) {
		_, msg, _ := conn.ReadMessage()
		println(string(msg))
	}(conn)
}

// Periodically send a JSON message to the client
func ws3(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, nil)
	go func(conn *websocket.Conn) {
		ch := time.Tick(5 * time.Second)
		for range ch {
			conn.WriteJSON(jsonMsg{
				UserName:  "damienstanton",
				FirstName: "Damien",
				LastName:  "Stanton",
			})
		}
	}(conn)
}

// Set websocket state.
// In this case, ws.CLOSED
func ws4(w http.ResponseWriter, r *http.Request) {
	var conn, _ = upgrader.Upgrade(w, r, nil)
	go func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
			}
		}
	}(conn)
}
