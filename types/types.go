package types

type Login struct {
	ClientId int    `json:"client_id"`
	Username string `json:"username"`
}

type WSMessage struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}

type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MessageType string

const (
	SignIn    MessageType = "login"
	PosUpdate             = "position_update"
)
