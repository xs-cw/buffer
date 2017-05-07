package buffer

import (
	"testing"
	"time"

	"github.com/wzshiming/ffmt"
	"github.com/wzshiming/task"
)

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
