package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const addr string = "localhost:8100"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWS(ch1 chan Query, ch2 chan Search) {
	fmt.Println("WS server running at:", addr)
	// http.HandleFunc("/", wsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// wsHandler(w, r, ch1, ch2)
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		go subWS(ws, ch1)
		go pubWS(ws, ch2)
	})

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func subWS(ws *websocket.Conn, ch chan<- Query) {
	for {
		_, message, err := ws.ReadMessage()

		if err != nil {
			log.Print("WS read:", err)
			break
		}

		query := Query{}
		json.Unmarshal(message, &query)

		fmt.Println("query:", query, query.Search)
		ch <- query
	}
}

func pubWS(ws *websocket.Conn, ch <-chan Search) {
	for {
		res := <-ch
		data, _ := json.Marshal(res)
		err := ws.WriteMessage(1, data)

		if err != nil {
			log.Println("WS write:", err)
			break
		}
	}
}
