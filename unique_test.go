// SPDX-License-Identifier: MIT

package unique

import (
	"context"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

func BenchmarkUnique(b *testing.B) {
	s := NewString()
	ctx, cancel := context.WithCancel(context.Background())
	go s.Serve(ctx)
	time.Sleep(time.Microsecond * 500)
	defer cancel()

	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	a.PanicString(func() {
		New(time.Microsecond, "", 10)
	}, "参数 duration 不能小于 1 秒")

	a.PanicString(func() {
		New(time.Second, "2006", 10)
	}, "无效的 prefixFormat 参数")

	a.PanicString(func() {
		New(time.Second, "20060102150405", 1)
	}, "参数 base 只能介于 [2,36] 之间")

	a.PanicString(func() {
		New(time.Second, "20060102150405", 37)
	}, "参数 base 只能介于 [2,36] 之间")

	u := New(time.Second, "20060102150405", 2)
	a.NotNil(u)
}

func TestIsValidDateFormat(t *testing.T) {
	a := assert.New(t, false)

	a.False(isValidDateFormat("2006"))
	a.False(isValidDateFormat("200601"))
	a.False(isValidDateFormat("20060102"))
	a.False(isValidDateFormat("2006010215"))
	a.False(isValidDateFormat("200601021504"))
	a.True(isValidDateFormat("20060102150405"))
	a.True(isValidDateFormat("05200601021504-"))
}

func TestUnique_String(t *testing.T) {
	a := assert.New(t, false)

	u := New(time.Second, "", 5)
	ctx, cancel := context.WithCancel(context.Background())
	go u.Serve(ctx)
	time.Sleep(time.Microsecond * 500)
	defer cancel()

	list := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 20)
		str := u.String()
		for _, item := range list {
			a.NotEqual(item, str)
		}
		list = append(list, str)
	}
}
