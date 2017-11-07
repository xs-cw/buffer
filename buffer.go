package buffer

import (
	"fmt"
	"time"

	"github.com/wzshiming/cache"
)

type Buffer struct {
	cache *cache.Memory
}

func NewBuffer() *Buffer {
	return &Buffer{
		cache: cache.NewMemory(),
	}
}

// Del 删除缓存数据
func (b *Buffer) Del(k string) {
	b.cache.Delete(k)
}

// Get 获取节点
func (b *Buffer) Get(k string, f MakeFunc) (*Node, bool) {
	i, ok := b.cache.GetOrPut(k, newNode(f), time.Second*10)
	n, _ := i.(*Node)
	return n, ok
}

// Buf 缓存数据
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {
	if f == nil {
		return nil, time.Time{}, fmt.Errorf("buff: 没有传入获取数据方法")
	}

	nn, ok := b.Get(k, f)

	i, t, e = nn.Latest()

	if !ok {
		if p := t.Sub(time.Now()); p > 0 {
			b.cache.SetTimeout(k, p)
		} else {
			b.cache.Delete(k)
		}
	}
	return

}
