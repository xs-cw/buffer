package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wzshiming/task"
)

type Buffer struct {
	Mut  sync.Mutex
	Buff map[string]*Node
	task *task.Task
}

func NewBuffer() *Buffer {
	return &Buffer{
		Buff: map[string]*Node{},
		task: task.NewTask(16),
	}
}

// Del 删除缓存数据
//  k:   缓存键
func (b *Buffer) Del(k string) {
	b.Mut.Lock()
	defer b.Mut.Unlock()
	delete(b.Buff, k)
}

// init 初始化缓存数据
//  k:   缓存键
func (b *Buffer) getNode(k string, f MakeFunc) *Node {
	b.Mut.Lock()
	defer b.Mut.Unlock()
	if b.Buff[k] == nil {
		t := NewNode(f)
		t.Update()
		b.Buff[k] = t
		if !t.timeout.IsZero() {
			b.task.Add(t.timeout, func() {
				b.Del(k)
			})
		}
		return t
	}
	return b.Buff[k]
}

// Buf 缓存数据
//  k:   缓存键
//  f:   刷新缓存的闭包
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("没有传入获取数据方法")
	}
	val := b.getNode(k, f)

	return val.Value()
}
