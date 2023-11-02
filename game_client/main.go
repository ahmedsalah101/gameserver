package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"gameserver/types"

	"github.com/gorilla/websocket"
)

var pl = fmt.Println

const wsServerEndpoint = "ws://localhost:4000/ws"

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func newGameClient(
	conn *websocket.Conn,
	username string,
) *GameClient {
	return &GameClient{
		clientID: rand.Intn(math.MaxInt),
		username: username,
		conn:     conn,
	}
}

func (c *GameClient) login() error {
	b, err := json.Marshal(types.Login{
		ClientID: c.clientID,
		Username: c.username,
	})
	if err != nil {
		return err
	}
	msg := types.WSMessage{
		Type: "Login",
		Data: b,
	}
	return c.conn.WriteJSON(msg)
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
	c := newGameClient(conn, "hello")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}
	// read loop
	go func() {
		var msg types.WSMessage
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("WS read error", err)
				continue
			}
			switch msg.Type {
			case "state":
				var state types.PlayerState
				if err := json.Unmarshal(msg.Data, &state); err != nil {
					fmt.Println("WS read error: ", err)
					continue
				}
				fmt.Println(
					"need to update the state of player",
					state,
				)
			default:
				fmt.Println(
					"receiving message we don't know",
				)
			}
			// fmt.Println("got message from the server", msg)
		}
	}()
	for {
		x := rand.Intn(1000)
		y := rand.Intn(1000)
		state := types.PlayerState{
			HP: 100,
			Postition: types.Postition{
				X: x,
				Y: y,
			},
		}
		b, err := json.Marshal(state)
		if err != nil {
			log.Fatal(err)
		}
		msg := types.WSMessage{
			Type: "PlayerState",
			Data: b,
		}
		if err := conn.WriteJSON(msg); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * 2000)
	}
}
