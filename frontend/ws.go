package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const addr string = "localhost:8100"

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws *websocket.Conn
)

func serveWS(ch1 chan Query, ch2 chan Search) {
	fmt.Println("WS server running at:", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		ws, err = upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		go subWS(ch1)
		go pubWS(ch2)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func subWS(ch chan<- Query) {
	for {
		_, message, err := ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("!!!error: %v", err)
			}
			break
		}

		query := Query{}
		json.Unmarshal(message, &query)

		fmt.Println("query:", query, query.Search)
		ch <- query
	}
}

func pubWS(ch <-chan Search) {
	for {
		res := <-ch
		data, _ := json.Marshal(res)
		err := ws.WriteMessage(1, data)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
