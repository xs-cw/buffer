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
func (b *Buffer) getNode(k string, f MakeFunc) (*Node, bool) {
	b.Mut.Lock()
	defer b.Mut.Unlock()
	if b.Buff[k] == nil {
		b.Buff[k] = NewNode(f)
		return b.Buff[k], true
	}
	return b.Buff[k], false
}

// Buf 缓存数据
//  k:   缓存键
//  f:   刷新缓存的闭包
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("没有传入获取数据方法")
	}
	val, bb := b.getNode(k, f)
	i, t, e = val.Value()
	if e != nil {
		return
	}
	if bb {
		b.task.Add(t, func() {
			b.Del(k)
		})
	}

	return
}
