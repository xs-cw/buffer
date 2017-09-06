package buffer

import (
	"fmt"
	"time"

	"github.com/wzshiming/cache"
)

type Buffer struct {
	buff cache.Cache
}

func NewBuffer() *Buffer {
	return &Buffer{
		buff: cache.NewMemory(),
	}
}

// Del 删除缓存数据
func (b *Buffer) Del(k string) {
	b.buff.Delete(k)
}

// Get 获取节点
func (b *Buffer) Get(k string) *Node {
	i := b.buff.Get(k)
	n, ok := i.(*Node)
	if ok {
		return n
	}
	return nil
}

// Buf 缓存数据
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("没有传入获取数据方法")
	}
	// 获取节点
	nn := b.Get(k)
	if nn != nil {
		return nn.Value()
	}

	nn = NewNode(f)
	i, t, e = nn.Latest()
	defer b.buff.Put(k, nn, t.Sub(time.Now()))
	return i, t, e
}
