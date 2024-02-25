package main

import (
	"encoding/json"
	"github.com/bulaioch/types"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
)

const wsServerEndpoint = "ws://localhost:40000/ws"

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func newGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		conn:     conn,
		clientID: rand.Intn(math.MaxInt),
		username: username,
	}
}

func (client *GameClient) login() error {
	login, err := json.Marshal(types.Login{
		ClientId: client.clientID,
		Username: client.username,
	})
	if err != nil {
		return err
	}
	msg := types.WSMessage{
		Type: "login",
		Data: login,
	}
	return client.conn.WriteJSON(msg)
}

func main() {

	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := newGameClient(conn, "Bob")
	if err := client.login(); err != nil {
		log.Fatal(err)
	}

}
