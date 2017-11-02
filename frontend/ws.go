package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WS) init(broker Broker) {
	http.HandleFunc("/", ws.upgrade)

	ws.chSearch = broker.chSearch
	ws.chReslt = broker.chReslt
}

func (ws *WS) serve() {
	fmt.Println("> WS server running at:", wsAddr)

	err := http.ListenAndServe(wsAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (ws *WS) upgrade(w http.ResponseWriter, r *http.Request) {
	var err error
	ws.conn, err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	go ws.sub()
	go ws.pub()
}

func (ws *WS) sub() {
	for {
		_, message, err := ws.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("%v", err)
			}
			break
		}

		query := Query{}

		json.Unmarshal(message, &query)
		fmt.Printf("[ws/sub] received search term \"%s\"\n", query.Search)

		ws.chSearch <- query
	}
}

func (ws *WS) pub() {
	for {
		res := <-ws.chReslt
		data, _ := json.Marshal(res)
		err := ws.conn.WriteMessage(1, data)

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("[ws/pub] pubbed search \"%s\" results \n", res.Term)
	}
}
