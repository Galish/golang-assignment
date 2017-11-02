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

func (ws *WS) init(chSearch chan SearchQuery, chResult chan SearchResult) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		ws.conn, err = upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		go ws.sub(chSearch)
		go ws.pub(chResult)
	})
}

func (ws *WS) serve() {
	fmt.Println("> WS server running at:", wsAddr)

	err := http.ListenAndServe(wsAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (ws *WS) sub(ch chan<- SearchQuery) {
	for {
		_, message, err := ws.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("%v", err)
			}
			break
		}

		query := SearchQuery{}

		json.Unmarshal(message, &query)
		fmt.Printf("[ws/sub] received search term \"%s\"\n", query.Term)

		// fmt.Printf("WS message: %+v\n", message)
		// fmt.Printf("WS query: %+v\n", query)
		// fmt.Println("WS query search:", query.Search)

		ch <- query
	}
}

func (ws *WS) pub(ch <-chan SearchResult) {
	for {
		res := <-ch
		data, _ := json.Marshal(res)
		err := ws.conn.WriteMessage(1, data)

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("[ws/pub] pubbed search \"%s\" results \n", res.Term)
	}
}
