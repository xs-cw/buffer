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
		task: task.NewTaskBuf(16, 1024),
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
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("没有传入获取数据方法")
	}
	val, bb := b.init(k, f)
	if !bb {
		i, t, e = val.Value()
		if t.After(time.Now()) {
			return i, t, e
		}
	}

	ok, e := val.Flash()
	if e != nil {
		return nil, time.Time{}, e
	}
	if !ok {
		i, t, e = val.Value()
		b.task.Add(t, func() {
			b.Del(k)
		})
	}
	return i, t, e
}
