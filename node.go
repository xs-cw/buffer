package buffer

import (
	"sync"
	"time"
)

type MakeFunc func() (interface{}, time.Time, error)

type Node struct {
	timeout time.Time   // 过期
	data    interface{} // 缓存数据
	err     error       // 错误
	fun     MakeFunc    // 更新缓存
	mut     sync.RWMutex
}

func newNode(f MakeFunc) *Node {
	n := &Node{
		timeout: time.Unix(0, 0),
		fun:     f,
	}
	return n
}

// IsValid 数据是有效的
func (n *Node) IsValid() bool {
	n.mut.RLock()
	defer n.mut.RUnlock()
	return n.timeout.After(time.Now())
}

// Latest 最新的缓存数据 返回 是否刷新数据
func (n *Node) Latest() (interface{}, time.Time, error) {
	n.mut.Lock()
	defer n.mut.Unlock()
	if n.timeout.After(time.Now()) {
		return n.data, n.timeout, n.err
	}
	n.data, n.timeout, n.err = n.fun()
	return n.data, n.timeout, n.err
}
