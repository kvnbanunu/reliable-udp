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

type SendSuccessMsg struct {
	SeqNum  uint8
	Retries uint8
	Timeout bool
}

type RecvSuccessMsg struct {
	Packet utils.Packet
}

type TimeoutMsg struct{}

type CancelMsg struct{}
