package types

type Login struct {
	ClientId int    `json:"client_id"`
	Username string `json:"username"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
