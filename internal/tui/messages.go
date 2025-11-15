package tui

import "reliable-udp/internal/utils"

type ErrMsg struct {
	err error
}

type LogMsg struct {
	MsgSent int `json:"messagesSent"`
	MsgRecv int `json:"messagesReceived"`
}

type SentMsg struct{}

type RecvMsg struct{
	Packet utils.Packet
}

type TimeoutMsg struct{}
