package types

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Login struct {
	ClientID int    `json:"clientID"`
	Username string `json:"username"`
}

type Postition struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type PlayerState struct {
	HP        int       `json:"hp"`
	Postition Postition `json:"postition"`
	SessionID int
	// Velocity int // may be to check anti-cheat
	// PowerUP int
}
