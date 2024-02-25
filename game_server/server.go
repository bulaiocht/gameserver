package main

import (
	"encoding/json"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/bulaioch/types"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"net/http"
)

// PlayerSession represents a player session with a unique session ID,
// indicating whether the player is in a lobby, and a reference to
// the player's connection.
type PlayerSession struct {
	sessionID   int
	inLobby     bool
	conn        *websocket.Conn
	playerState types.PlayerState
}

func (session *PlayerSession) Receive(context *actor.Context) {
	switch context.Message().(type) {
	case actor.Started:
		session.readLoop()
	}
}

func (session *PlayerSession) handleWSMessage(msg types.WSMessage) {
	switch msg.Type {
	case types.SignIn:
		var login types.Login
		if err := json.Unmarshal(msg.Data, &login); err != nil {
			log.Printf("Error unmarshalling login data: %s", err)
			panic(err)
		}
		if e := session.handleLogin(login); e != nil {
			log.Printf("Unable to login client: %s", e)
			panic(e)
		}
	case types.PosUpdate:
		var pos types.Position
		if err := json.Unmarshal(msg.Data, &pos); err != nil {
			log.Printf("Error unmarshalling login data: %s", err)
			panic(err)
		}
		log.Printf("position update: x=%d, y=%d", pos.X, pos.Y)
		session.playerState.Position.X = pos.X
		session.playerState.Position.Y = pos.Y

	}
}

func (session *PlayerSession) handleLogin(msg types.Login) error {
	log.Printf("Received login message: %+v", msg)
	return nil
}

func (session *PlayerSession) readLoop() {
	var msg types.WSMessage
	for {
		err := session.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading incomming message %s:", err)
			return
		}
		go session.handleWSMessage(msg)
	}
}

func newPlayerSession(sessionID int, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionID: sessionID,
			conn:      conn,
			playerState: types.PlayerState{
				Health: 100,
				Position: types.Position{
					X: 0,
					Y: 0,
				},
			},
		}
	}
}

// GameServer represents a game server that handles WebSocket connections and game sessions.
type GameServer struct {
	ctx      *actor.Context
	sessions map[*actor.PID]struct{}
}

func newGameServer() actor.Receiver {
	return &GameServer{
		sessions: make(map[*actor.PID]struct{}),
	}
}

func (server *GameServer) Receive(context *actor.Context) {
	switch msg := context.Message().(type) {
	case actor.Initialized:
		log.Print("TODO")
	case actor.Started:
		server.ctx = context
		server.startHTTP()
		_ = msg
	}
}

func (server *GameServer) startHTTP() {
	log.Print("starting an HTTP server on port 40000")
	go func() {
		http.HandleFunc("/ws", server.handleWS)
		_ = http.ListenAndServe(":40000", nil)
	}()
}

func (server *GameServer) handleWS(writer http.ResponseWriter, req *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(writer, req, req.Header)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	log.Print("Client is trying to connect")
	log.Print(conn)
	sid := rand.Intn(math.MaxInt)
	session := newPlayerSession(sid, conn)
	pid := server.ctx.SpawnChild(session, fmt.Sprintf("session_%d", sid))
	server.sessions[pid] = struct{}{}
	log.Printf("client with sid: %d and pid: %s just connected", sid, pid)
}

// main is the entry point of the program.
// It initializes an engine config, creates an engine and spawns a new game server actor.
// After that, it enters an endless loop.
func main() {
	config := actor.NewEngineConfig()
	engine, _ := actor.NewEngine(config)
	log.Print("Spawning an actor")
	engine.Spawn(newGameServer, "server")
	//this is endless loop
	select {}
}
