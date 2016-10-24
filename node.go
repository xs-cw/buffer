package buffer

import (
	"time"

	"sync"
)

type MakeFunc func() (interface{}, time.Time, error)

var node = struct{}{}

type Node struct {
	Timeout time.Time   // 过期
	Data    interface{} // 缓存数据
	Func    MakeFunc    // 更新缓存
	mut     sync.Mutex
}

func NewNode() *Node {
	return &Node{
		Timeout: time.Unix(0, 0),
	}
}

// 是有效的
func (n *Node) IsValid() bool {
	n.mut.Lock()
	defer n.mut.Unlock()
	return n.Timeout.After(time.Now())
}

// Value 数据
func (n *Node) Value() (interface{}, time.Time, error) {
	n.mut.Lock()
	defer n.mut.Unlock()
	return n.Data, n.Timeout, nil
}

// Update 更新数据
func (n *Node) Update() error {
	n.mut.Lock()
	defer n.mut.Unlock()

	d, t, err := n.Func()
	if err == nil {
		n.Data = d
		n.Timeout = t
	}
	return err
}

// Flash 刷新
func (n *Node) Flash() (bool, error) {
	if !n.IsValid() {
		err := n.Update()
		if err != nil {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
