package common

type MessageType string

const (
	SystemInfoUpdateMsgType MessageType = "system_info_update"
	RegisterMsgType         MessageType = "register"
	SessionAckMsgType       MessageType = "session_ack"
	SessionRemoteExecReq    MessageType = "session_remote_exec_request"
	SessionRemoteExecResp   MessageType = "session_remote_exec_response"
	SessionRemoteExecInit   MessageType = "session_remote_exec_init"
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

type RegisterMsg struct {
	Action         string     `json:"action"` // should be "register"
	ClientUsername string     `json:"client_username"`
	JoineeType     JoineeType `json:"joinee_type"` // choices{worker,consumer}
}

type SessionAckMsg struct {
	SessionId string `json:"session_id"`
}

type SessionRemoteExecInitMsg struct {
	ExecutionID string `json:"execution_id"`
}
