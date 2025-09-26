package common

type MessageType string

const (
	SystemInfoUpdateMsgType MessageType = "system_info_update"
	RegisterMsgType         MessageType = "register"
	SessionAckMsgType       MessageType = "session_ack"
	SessionExecReqMsgType   MessageType = "session_exec_request"
)

type Message struct {
	Type        MessageType `json:"type"`
	Description string      `json:"description"`
	Body        any         `json:"body"`
}

type JoineeType string

const (
	Worker   JoineeType = "worker"
	Consumer JoineeType = "consumer"
)

type RegisterMessage struct {
	Action         string     `json:"action"` // should be "register"
	ClientUsername string     `json:"client_username"`
	JoineeType     JoineeType `json:"joinee_type"` // choices{worker,consumer}
}

type SessionAckMessage struct {
	SessionId string `json:"session_id"`
}
