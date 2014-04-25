package doracle

import (
	"sync"
)

type Oracle struct {
	maxTs int64
	mutex sync.Mutex
}

func NewOracle() *Oracle {
	return &Oracle{
		maxTs: -1,
	}
}

func (o *Oracle) GetTimestamp(num int32) int64 {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.maxTs += int64(num)
	return o.maxTs
}
