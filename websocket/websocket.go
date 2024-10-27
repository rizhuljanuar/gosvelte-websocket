package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type WSServer struct {
	connections		map[string]map[*websocket.Conn]struct{}
	connectionMutex sync.Mutex
}

func NewWebServer() *WSServer {
	return &WSServer {
		connections: make(map[string]map[*websocket.Conn]struct{}),
	}
}

type WsResponseData struct {
	// Action string `json:"action"`
	Message string `json:"message"`
	// SenderId string `json:"senderId"`
	// Name string `json:"name"`
	// CreatedAt string `json:"CreatedAt"`
}

func (w *WSServer) Broadcast(chatRoomId, msg string) {
	w.connectionMutex.Lock()
	defer w.connectionMutex.Unlock()

	if groupConnections, ok := w.connections[chatRoomId]; ok {
		for conn := range groupConnections {
			if err := websocket.Message.Send(conn, msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func (w *WSServer) WsHandler(c echo.Context) error {
	chatRoomId := c.Param("chat_room_id")

	websocket.Handler(func (ws *websocket.Conn) {
		defer ws.Close()

		log.Println("Connection: ", chatRoomId)

		w.connectionMutex.Lock()
		if _, ok := w.connections[chatRoomId]; !ok {
			w.connections[chatRoomId] = make(map[*websocket.Conn]struct{})
		}

		w.connections[chatRoomId][ws] = struct{}{}
		w.connectionMutex.Unlock()

		defer func() {
			w.connectionMutex.Lock()
			delete(w.connections[chatRoomId], ws)
			if (len(w.connections[chatRoomId]) == 0) {
				delete(w.connections, chatRoomId)
			}

			w.connectionMutex.Unlock()

			log.Println("Disconnected: ", ws.RemoteAddr(), "from: ", chatRoomId)
		}()

		for {
			var wsres = WsResponseData{}

			if err := websocket.JSON.Receive(ws, &wsres); err != nil {
				log.Println("Err when receiving data: ", err)
				break
			}

			jsonRes, err := json.Marshal(wsres)

			if err != nil {
				log.Println("Error when mashaling response: ", err)
			}

			log.Println("Incoming message: ", string(jsonRes), "for: ", chatRoomId)

			w.Broadcast(chatRoomId, string(jsonRes))
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
