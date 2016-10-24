package buffer

import (
	"fmt"
	"testing"
	"time"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/fork"
)

func TestA(t *testing.T) {
	buf := NewBuffer()
	fo := fork.NewFork(100)
	i := 0
	for ; i != 100; i++ {
		fo.Push(func() {
			b, _, err := buf.Buf("hello", func() (interface{}, time.Time, error) {
				return "world " + fmt.Sprint(time.Now()), time.Now().Add(time.Second / 2), nil
			})
			ffmt.Mark(b, err)
		})
	}
	ffmt.Mark(i)
	<-time.After(time.Second)
	ffmt.Mark(i)
	for ; i != 200; i++ {
		fo.Push(func() {
			b, _, err := buf.Buf("hello", func() (interface{}, time.Time, error) {
				return "world " + fmt.Sprint(time.Now()), time.Now().Add(time.Second / 2), nil
			})
			ffmt.Mark(b, err)
		})
	}
	ffmt.Mark(i)
	<-time.After(time.Second)
	//	for {
	//		//ffmt.Puts(buf)
	//		runtime.Gosched()
	//	}
}
