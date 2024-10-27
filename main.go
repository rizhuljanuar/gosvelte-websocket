package main

import (
	"fmt"
	"gosvelt-websocket/websocket"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Hello world")

	e := echo.New()
	wsServer := websocket.NewWebServer()

	e.GET("/cek", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"Hello":"Word"})
	})

	e.GET("/ws/:chat_room_id", wsServer.WsHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
