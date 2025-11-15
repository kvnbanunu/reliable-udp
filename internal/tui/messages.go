package tui

import "reliable-udp/internal/utils"

type ErrMsg struct {
	Err error
}

type LogMsg struct {
	MsgSent int `json:"messagesSent"`
	MsgRecv int `json:"messagesReceived"`
}

type LogSuccessMsg struct{}

type SendSuccessMsg struct{}

type RecvSuccessMsg struct {
	Packet utils.Packet
}

type TimeoutMsg struct{}
