package PBFT

import "time"

type Message struct {
	From      ID
	To        ID
	TimeStamp time.Time
	//目前有pre-prepare、prepare、commit
	MsgType string
	MsgBody []byte
}

// NewMessage 创建新消息
func NewMessage(msgType string, msgBody []byte) *Message {
	return &Message{
		From:      ID{},
		To:        ID{},
		TimeStamp: time.Time{},
		MsgType:   msgType,
		MsgBody:   msgBody,
	}
}
