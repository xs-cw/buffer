package buffer

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/task"
)

func makeRand() string {
	s := make([]byte, 64)
	rand.Read(s)
	return base64.StdEncoding.EncodeToString(s)
}

var buf = NewBuffer()

func benchmarkA(b *testing.B, c, d time.Duration) {
	for i := 0; i < b.N; i++ {
		buf.Buf(b.Name(), func() (interface{}, time.Time, error) {
			rr := makeRand()
			//ffmt.Mark("flash", b.Name(), time.Now().Format(time.RFC3339Nano), rr)
			if c > 0 {
				time.Sleep(c)
			}
			return rr, time.Now().Add(d), nil
		})
	}
	return
}

// 测试默认1s缓存情况下效率
func BenchmarkAZS(b *testing.B) {
	benchmarkA(b, 0, time.Second)
}

func BenchmarkAS10S(b *testing.B) {
	benchmarkA(b, time.Second/10, time.Second)
}

func BenchmarkAS2S(b *testing.B) {
	benchmarkA(b, time.Second/2, time.Second)
}

func BenchmarkB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeRand()
	}
}

func BenchmarkB2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf.Buf(b.Name(), func() (interface{}, time.Time, error) {
			return makeRand(), time.Now(), nil
		})
	}
}

func BenchmarkP(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkA(b, time.Second/100, time.Second)
		}
	})
}

func BenchmarkPS10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkA(b, time.Second/10, time.Second)
		}
	})
}

func TestA(t *testing.T) {
	buf := NewBuffer()
	fo := task.NewTask(100)
	i := 0
	uc := 0
	pc := 0
	for ; i != 100; i++ {
		cc := time.Second / 50 * (time.Duration(i) / 2)
		fo.Add(time.Now().Add(cc), func() {
			b, t, err := buf.Buf("hello", func() (interface{}, time.Time, error) {
				v, t := time.Now(), time.Now().Add(time.Second/10)
				ffmt.Mark("update key", v, t)
				uc++
				//time.Sleep(time.Second / 100)
				return v, t, nil
			})
			pc++
			now := time.Now()
			ffmt.Mark(now.Sub(b.(time.Time)), t.Sub(now), err)
		})
	}

	fo.Join()
	ffmt.Mark(uc, pc)
}
