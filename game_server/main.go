package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"

	"gameserver/types"
)

type PlayerSession struct {
	ctx       *actor.Context
	username  string
	sessionID int
	clientID  int
	inLobby   bool
	conn      *websocket.Conn
	serverPID *actor.PID
}

func newPlayerSession(serverPID *actor.PID, sid int,
	conn *websocket.Conn,
) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionID: sid,
			conn:      conn,
			serverPID: serverPID,
		}
	}
}

func (ps *PlayerSession) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		ps.ctx = c
		go ps.readLoop()
	case *types.PlayerState:
		ps.sendPlayerState(msg)
	default:
		fmt.Println("rev: ", msg)
	}
}

func (s *PlayerSession) sendPlayerState(
	state *types.PlayerState,
) {
	b, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	msg := types.WSMessage{
		Type: "state",
		Data: b,
	}

	if err := s.conn.WriteJSON(msg); err != nil {
		panic(err)
	}
	fmt.Println("sending this state ", state)
}

func (s *PlayerSession) readLoop() {
	var msg types.WSMessage
	for {
		if err := s.conn.ReadJSON(&msg); err != nil {
			fmt.Println("read error", err)
			return
		}
		go s.handleMessage(msg)
	}
}

func (s *PlayerSession) handleMessage(msg types.WSMessage) {
	switch msg.Type {
	case "Login":
		fmt.Println("recieved login msg")
		var loginMsg types.Login
		if err := json.Unmarshal(msg.Data, &loginMsg); err != nil {
			panic(err)
		}
		// TODO: auth
		s.clientID = loginMsg.ClientID
		s.username = loginMsg.Username
	case "PlayerState":
		var ps types.PlayerState
		if err := json.Unmarshal(msg.Data, &ps); err != nil {
			panic(err)
		}
		ps.SessionID = s.sessionID
		if s.ctx != nil {
			s.ctx.Send(s.serverPID, &ps)
		}
	default:
		fmt.Println("recieved unkown message: ", msg)
	}
}

type GameServer struct {
	ctx      *actor.Context
	sessions map[int]*actor.PID
}

func newGameServer() actor.Receiver {
	return &GameServer{
		sessions: make(map[int]*actor.PID),
	}
}

func (s *GameServer) bcast(
	from *actor.PID,
	state *types.PlayerState,
) {
	for _, pid := range s.sessions {
		if !pid.Equals(from) {
			s.ctx.Send(pid, state)
		}
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case *types.PlayerState:
		s.bcast(c.Sender(), msg)
	case actor.Started:
		s.startHTTP()
		s.ctx = c
		_ = msg
	}
}

func (s *GameServer) startHTTP() {
	fmt.Println("starting HTTP server on port 4000")
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(":4000", nil)
	}()
}

// const WS_ENDPOINT = ""

// handles the upgrade of websocket
func (s *GameServer) handleWS(
	w http.ResponseWriter,
	r *http.Request,
) {
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	// Deprecated
	// conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println("ws upgrade error: ", err)
		return
		// panic(err)
	}

	fmt.Println("new client trying to connect")
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(
		newPlayerSession(s.ctx.PID(), sid, conn),
		fmt.Sprintf("session_%d", sid),
	)
	s.sessions[sid] = pid
	fmt.Printf(
		"client with sid %d and pid %s just connected\n",
		sid,
		pid,
	)
	conn.WriteJSON(map[string]string{"msg": "success"})
}

func main() {
	e := actor.NewEngine()
	e.Spawn(newGameServer, "server")
	fmt.Println("server !")
	select {}
}

type HTTPServer struct{}
