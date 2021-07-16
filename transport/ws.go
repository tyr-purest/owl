package transport

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// webSocket -
type webSocket struct {
	mux   sync.RWMutex
	conn  *websocket.Conn
	out   chan []byte
	Close chan bool
}

// NewWs -
func NewWs(w http.ResponseWriter, r *http.Request) (*webSocket, error) {
	subprotocols := r.Header["Sec-Websocket-Protocol"]
	upgrader.Subprotocols = subprotocols

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("An error occured while upgrading the connection: %v\n", err)
		return nil, err
	}

	ws := &webSocket{
		conn:  conn,
		out:   make(chan []byte),
		Close: make(chan bool),
	}

	go ws.reader()
	go ws.writer()

	return ws, nil
}

// Reader -
func (ws *webSocket) reader() {
	defer ws.connClose()

	for {
		_, _, err := ws.conn.ReadMessage()

		// reader err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("webSocket error")
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("Client reload disconnected")
				return
			}

			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				log.Println("Client actively disconnected")
				return
			}

			log.Println("other error")

			break
		}

		//spew.Dump(message)
	}
}

// Writer -
func (ws *webSocket) writer() {
	defer ws.connClose()
	for {
		select {
		case message, ok := <-ws.out:
			if !ok {
				ws.conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
			}

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			ws.mux.Lock()
			w.Write(message)
			w.Close()
			ws.mux.Unlock()
		}
	}
}

// connClose -
func (ws *webSocket) connClose() {
	ws.Close <- true
	ws.conn.Close()
}
