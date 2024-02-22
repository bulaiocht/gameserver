package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type PlayerSession struct {
	clientID int
	username string
	inLobby  bool
	conn     *websocket.Conn
}

func (session *PlayerSession) Receive(context *actor.Context) {
	switch msg := context.Message().(type) {
	case actor.Initialized:
		log.Print("Player session initialized")
	case actor.Started:
		log.Print("Player session Started")
		log.Printf("Message received: %s", msg)
		_ = msg
	}
}

func newPlayerSession(clientID int, username string, conn *websocket.Conn) actor.Receiver {
	return &PlayerSession{
		clientID: clientID,
		username: username,
		conn:     conn,
	}
}

type GameServer struct{}

func newGameServer() actor.Receiver {
	return &GameServer{}
}

func (server *GameServer) startHTTP() {
	log.Print("starting an HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", server.handleWS)
		_ = http.ListenAndServe(":40000", nil)
	}()
}

func (server *GameServer) Receive(context *actor.Context) {
	switch msg := context.Message().(type) {
	case actor.Initialized:
		log.Print("TODO")
	case actor.Started:
		server.startHTTP()
		_ = msg
	}
}

func (server *GameServer) handleWS(writer http.ResponseWriter, req *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(writer, req, req.Header)
	if err != nil {
		log.Printf("upgrade error: %s", err)
	}
	log.Print("Client is trying to connect")
	log.Print(conn)
}

func main() {
	config := actor.NewEngineConfig()
	engine, _ := actor.NewEngine(config)
	log.Print("Spawning an actor")
	engine.Spawn(newGameServer, "server")
	select {}
}
