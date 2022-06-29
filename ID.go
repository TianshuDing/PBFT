package PBFT

import (
	"fmt"
	"strconv"
)

type ID struct {
	ShardID      int
	IntraShardID int
	IPAddr       string
}

// ToString 返回一个字符串类型数据，表示节点全局唯一ID
// 格式为分片ID-片内ID
func (i *ID) ToString() string {
	v1 := strconv.FormatInt(int64(i.ShardID), 10)
	v2 := strconv.FormatInt(int64(i.IntraShardID), 10)
	return fmt.Sprintf("%s-%s", v1, v2)
}
