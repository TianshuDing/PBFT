package PBFT

// PhaseCount 阶段投票计数器
type PhaseCount struct {
	//累计计数器
	AckCount int
	//最后一次投票的Term Num
	//当出现高于当前Term Num的有效投票时
	//重置AckCount并更新LatestTerm为最新的Term Num
	LatestTerm int
	//消息管道
	MsgCh chan *Message
}

// NewPhaseCount 创建新的PhaseCount实例，输入一个节点总数参数
// 该参数用于生成对应的消息缓存管道长度设置，默认为3倍节点总数
func NewPhaseCount(nodeNum int) *PhaseCount {
	value := new(PhaseCount)
	value.AckCount = 0
	value.LatestTerm = 0
	value.MsgCh = make(chan *Message, 3*nodeNum)

	return value
}
