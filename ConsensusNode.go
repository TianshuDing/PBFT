package PBFT

import (
	"crypto/rand"
	"fmt"
	"go.uber.org/zap"
	"math"
	"math/big"
)

// ConsensusNode 共识节点结构体
type ConsensusNode struct {
	//日志打印器
	Logger *zap.Logger
	//节点角色，Primary、follower、seed
	Role string
	//视图编号，主节点无故障则保持不变
	ViewNum int
	//共识提案编号，每个提案（区块）的独立编号
	ProposalNum int
	//节点接收足够Prepare
	PreparePhaseCount   *PhaseCount
	PrepareCheckStartCh chan struct{}
	//节点接收足够Commit
	CommitPhaseCount   *PhaseCount
	CommitCheckStartCh chan struct{}
	//节点结束通知管道
	DoneCh chan struct{}
	//指令系统
	CLI chan string
	//内存模拟的网络
	NetWork *MemoryNetWork

	/***********************/
	//seedNode特有数据结构
	CLIs map[ID]chan string
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		//no error
	}
}

// Run 共识节点主循环
func (csn *ConsensusNode) Run() {

	for {
		select {
		case <-csn.DoneCh:
			return
		case msg := <-csn.NetWork.LocalMsgCh:
			csn.HandleMsg(msg)
		case cli := <-csn.CLI:
			csn.HandleCLI(cli)
		case <-csn.PrepareCheckStartCh:
			go csn.PrepareMsgCheck()
		case <-csn.CommitCheckStartCh:
			go csn.CommitMsgCheck()
		}
	}

}

// HandleMsg 消息处理模块，根据不同的交易类型进行分类
// Warning: some msg processing is time-consuming
// please using goroutine
func (csn *ConsensusNode) HandleMsg(msg *Message) {
	switch msg.MsgType {

	case "Pre-Prepare":
		info := fmt.Sprintf("[%s] sending Prepare Msg", csn.NetWork.NodeID.ToString())
		csn.Logger.Debug(info)
		csn.NetWork.BroadCast("all", 0, NewMessage("Prepare", nil))
		//收到pre-prepare消息的副本节点开始检测prepare消息
		csn.PrepareCheckStartCh <- struct{}{}

	case "Prepare":
		csn.PreparePhaseCount.MsgCh <- msg
	case "Commit":
		csn.CommitPhaseCount.MsgCh <- msg
	}
}

// HandleCLI 指令处理响应机制
func (csn *ConsensusNode) HandleCLI(cli string) {
	switch cli {
	case "/hello":
		info := fmt.Sprintf("[%s] Hello World! Time to fly", csn.NetWork.NodeID.ToString())
		csn.Logger.Info(info)
	case "/start":
		if csn.Role == "seed" {
			/**********************************************************************/
			//随机选择一个节点担任主节点
			randomNode, _ := rand.Int(rand.Reader, big.NewInt(int64(len(csn.CLIs))))
			i := int64(0)
			for _, v := range csn.CLIs {
				if i == randomNode.Int64() {
					v <- "/pre-prepare"
					break
				} else {
					i++
				}
			}
			/**********************************************************************/
		} else {
			//you are not seed node
		}
	case "/pre-prepare":
		info := fmt.Sprintf("[%s] sending Pre-Prepare Msg", csn.NetWork.NodeID.ToString())
		csn.Logger.Debug(info)

		csn.NetWork.BroadCast("all", 0, NewMessage("Pre-Prepare", nil))
		//发送pre-prepare消息的主节点开始检测prepare消息
		csn.PrepareCheckStartCh <- struct{}{}
	}

}

// PrepareMsgCheck Prepare阶段
func (csn *ConsensusNode) PrepareMsgCheck() {
	//计算拜占庭容错数量上限
	Q1 := int(math.Floor(float64(len(csn.NetWork.RemoteNodes)) / 3))
	if Q1 == 0 {
		Q1 = 1
	} else {
		//do nothing
	}
	for {
		select {
		case <-csn.PreparePhaseCount.MsgCh:
			csn.PreparePhaseCount.AckCount++
			//收到超过2f个prepare消息
			if csn.PreparePhaseCount.AckCount == 2*Q1 {
				info := fmt.Sprintf("[%s] sending Commit Msg", csn.NetWork.NodeID.ToString())
				csn.Logger.Debug(info)
				csn.NetWork.BroadCast("all", 0, NewMessage("Commit", nil))
				return
			} else {
				continue
			}
		}
	}
}

// CommitMsgCheck Commit阶段
func (csn *ConsensusNode) CommitMsgCheck() {
	//计算拜占庭容错数量上限
	Q2 := int(math.Floor(float64(len(csn.NetWork.RemoteNodes)) / 3))
	if Q2 == 0 {
		Q2 = 1
	} else {
		//do nothing
	}
	for {
		select {
		case <-csn.CommitPhaseCount.MsgCh:
			csn.CommitPhaseCount.AckCount++
			//收到超过2f个commit消息
			if csn.CommitPhaseCount.AckCount == 2*Q2 {
				info := fmt.Sprintf("[%s] get enough Commit Msg", csn.NetWork.NodeID.ToString())
				csn.Logger.Debug(info)

				return
			} else {
				continue
			}
		}
	}
}
