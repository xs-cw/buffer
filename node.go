package buffer

import (
	"time"

	"sync"
)

type MakeFunc func() (interface{}, time.Time, error)

type Node struct {
	timeout time.Time   // 过期
	data    interface{} // 缓存数据
	fun     MakeFunc    // 更新缓存
	mut     sync.RWMutex
}

func NewNode(f MakeFunc) *Node {
	n := &Node{
		timeout: time.Unix(0, 0),
		fun:     f,
	}
	n.Update()
	return n
}

// IsValid 数据是有效的
func (n *Node) IsValid() bool {
	n.mut.RLock()
	defer n.mut.RUnlock()
	return n.timeout.After(time.Now())
}

// Latest 最新的缓存数据
func (n *Node) Latest() (interface{}, time.Time, error) {
	if n.IsValid() {
		return n.Value()
	}
	err := n.Update()
	if err != nil {
		return nil, time.Time{}, err
	}
	return n.Value()
}

// Value 数据
func (n *Node) Value() (interface{}, time.Time, error) {
	n.mut.RLock()
	defer n.mut.RUnlock()
	return n.data, n.timeout, nil
}

// Update 更新数据
func (n *Node) Update() error {
	n.mut.Lock()
	defer n.mut.Unlock()
	d, t, err := n.fun()
	if err == nil {
		n.data = d
		n.timeout = t
	}
	return err
}

// Flash 刷新
func (n *Node) Flash() (bool, error) {
	if !n.IsValid() {
		err := n.Update()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
