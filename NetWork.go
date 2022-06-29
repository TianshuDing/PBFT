package PBFT

import (
	"errors"
	"fmt"
	"time"
)

// NetWork 网络接口，拥有发送和接收消息两个方法
type NetWork interface {
	Send(msg *Message, to ID) error
	Receive()
}

// MemoryNetWork 内存共享模拟网络
type MemoryNetWork struct {
	//共识节点ID
	NodeID      ID
	LocalMsgCh  chan *Message
	RemoteNodes map[ID]chan *Message
	ShardNodes  map[int]map[ID]chan *Message
}

// Send 发送消息方法，直接将消息发送至节点的接收消息管道
// 若节点不存在，则返回error；否则为nil
func (mnw *MemoryNetWork) Send(msg *Message, to ID) error {

	/************************/
	//新建消息时只配置了消息类型和消息体
	//在发送消息最后一步配置时间戳和发送接收双方ID
	msg.From = mnw.NodeID
	msg.To = to
	msg.TimeStamp = time.Now()
	/************************/

	if value, ok := mnw.RemoteNodes[to]; ok {
		value <- msg
	} else {
		errInfo := fmt.Sprintf("%s don't exist", to.ToString())
		return errors.New(errInfo)
	}

	return nil
}

func (mnw *MemoryNetWork) Receive() {

}

// BroadCast 消息广播方法
// scale表示范围（all代表所有节点，shard代表某一分片），shard代表分片编号（只有scale为shard时才有效）
func (mnw *MemoryNetWork) BroadCast(scale string, shard int, msg *Message) {

	switch scale {
	case "all":
		for k1, v1 := range mnw.RemoteNodes {
			if k1.ToString() == mnw.NodeID.ToString() {
				//do not send to myself
			} else {
				//warning: remote node msg receive channel may be full
				//this simple method may cause problem

				/************************/
				//新建消息时只配置了消息类型和消息体
				//在发送消息最后一步配置时间戳和发送接收双方ID
				msg.From = mnw.NodeID
				msg.To = k1
				msg.TimeStamp = time.Now()
				/************************/

				v1 <- msg
			}
		}
	case "shard":
		//warning: the shard may not exist
		for k2, v2 := range mnw.ShardNodes[shard] {
			if k2.ToString() == mnw.NodeID.ToString() {
				//do not send to myself
			} else {
				//warning: remote node msg receive channel may be full
				//this simple method may cause problem

				/************************/
				//新建消息时只配置了消息类型和消息体
				//在发送消息最后一步配置时间戳和发送接收双方ID
				msg.From = mnw.NodeID
				msg.To = k2
				msg.TimeStamp = time.Now()
				/************************/

				v2 <- msg
			}
		}
	}

}
