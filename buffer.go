package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wzshiming/cache"
)

type Buffer struct {
	cache *cache.Memory
	mut   sync.Mutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		cache: cache.NewMemory(),
	}
}

// Del 删除缓存数据
func (b *Buffer) del(k string) {
	b.cache.Delete(k)
}

// Get 获取节点
func (b *Buffer) get(k string) *Node {
	i := b.cache.Get(k)
	n, ok := i.(*Node)
	if ok {
		return n
	}
	return nil
}

// Buf 缓存数据
func (b *Buffer) Buf(k string, f MakeFunc) (i interface{}, t time.Time, e error) {

	if f == nil {
		return nil, time.Time{}, fmt.Errorf("buff: 没有传入获取数据方法")
	}

	nn := b.get(k)
	// 获取节点
	if nn != nil {
		return nn.Latest()
	}

	// 加锁加载
	b.mut.Lock()
	defer b.mut.Unlock()
	if nn = b.get(k); nn == nil {
		nn = newNode(f)
		b.cache.Put(k, nn, 0)
	}

	i, t, e = nn.Latest()

	if p := t.Sub(time.Now()); p > 0 {
		b.cache.SetTimeout(k, p)
	} else {
		b.cache.Delete(k)
	}

	return i, t, e

}
