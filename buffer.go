package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wzshiming/fork"
)

type Buffer struct {
	Mut  sync.Mutex
	Buff map[string]*Node
	fork *fork.Fork
}

func NewBuffer() *Buffer {
	return &Buffer{
		Buff: map[string]*Node{},
		fork: fork.NewFork(2),
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
func (b *Buffer) init(k string, f MakeFunc) (*Node, bool) {
	b.Mut.Lock()
	defer b.Mut.Unlock()
	bb := false
	if b.Buff[k] == nil {
		bb = true
		b.Buff[k] = NewNode()
	}
	b.Buff[k].Func = f
	return b.Buff[k], bb
}

// Buf 缓存数据
//  k:   缓存键
//  f:   缓存执行的动作
func (b *Buffer) Buf(k string, f MakeFunc) (interface{}, time.Time, error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("没有传入获取数据方法")
	}
	val, _ := b.init(k, f)

	i, t, e := val.Value()
	if t.Before(time.Now()) {
		b.fork.Push(func() {
			val.Flash()
		})
		b.fork.Join()
		if i == nil {
			i, t, e = val.Value()
		}
	}
	return i, t, e
}
