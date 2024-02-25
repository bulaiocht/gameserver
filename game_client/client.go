package main

import (
	"encoding/json"
	"fmt"
	"github.com/bulaioch/types"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"log"
	"math"
	"math/rand"
	"time"
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
		Type: types.SignIn,
		Data: login,
	}
	return client.conn.WriteJSON(msg)
}

func (client *GameClient) writeLoop() error {
	for {
		time.Sleep(time.Millisecond * 300)
		if e := client.positionUpdate(); e != nil {
			log.Printf("Unable to update player position: %s", e)
			return e
		}
	}
}

func (client *GameClient) positionUpdate() error {

	x := rand.Intn(1000)
	y := rand.Intn(1000)
	position := types.Position{
		X: x,
		Y: y,
	}
	log.Printf("sending position update: %s", fmt.Sprint(position))

	posData, e := json.Marshal(position)
	if e != nil {
		log.Printf("Error, converting message: %s", e)
		return e
	}

	message := types.WSMessage{
		Type: types.PosUpdate,
		Data: posData,
	}

	if err := client.conn.WriteJSON(message); err != nil {
		log.Printf("Error sending message: %s", err)
		return err
	}

	return nil
}

func main() {

	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	conn, _, err := dialer.DialContext(timeout, wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := newGameClient(conn, "Bob")
	if err := client.login(); err != nil {
		log.Fatal(err)
	}

	if e := client.writeLoop(); e != nil {
		log.Fatal(e)
	}

}
