// SPDX-License-Identifier: MIT

package unique

import (
	"math"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/autoinc"
)

func BenchmarkUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String().String()
	}
}

func TestNew(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() {
		New(20, 0, 100, "", 10)
	})

	a.Panic(func() {
		New(20, 1, 0, "", 10)
	})

	a.Panic(func() {
		New(20, 1, 1, "2006", 10)
	})

	a.Panic(func() {
		New(20, 1, 1, "20060102150405", 1)
	})

	a.Panic(func() {
		New(20, 1, 1, "20060102150405", 37)
	})

	a.NotPanic(func() {
		New(20, 1, time.Second, "20060102150405", 2)
	})
}

func TestNumber(t *testing.T) {
	a := assert.New(t)

	n1 := NewNumber()
	n2 := NewNumber()
	a.NotEqual(n1, n2)

	n1 = Number()
	n2 = Number()
	a.Equal(n1, n2)
}

func TestIsValidDateFormat(t *testing.T) {
	a := assert.New(t)

	a.False(isValidDateFormat("2006"))
	a.False(isValidDateFormat("200601"))
	a.False(isValidDateFormat("20060102"))
	a.False(isValidDateFormat("2006010215"))
	a.False(isValidDateFormat("200601021504"))
	a.True(isValidDateFormat("20060102150405"))
	a.True(isValidDateFormat("05200601021504-"))
}

func TestUnique_String(t *testing.T) {
	a := assert.New(t)

	u := New(time.Now().Unix(), 100000, 50*time.Second, "", 5)
	list := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		str := u.String()
		for _, item := range list {
			a.NotEqual(item, str)
		}
		list = append(list, str)
	}
}

func TestUnique_String_overflow(t *testing.T) {
	a := assert.New(t)

	u := New(time.Now().Unix(), 100000, 50*time.Minute, "", 5)
	u.ai = autoinc.New(math.MaxInt64-1, 2, 2)
	list := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		str := u.String()
		for _, item := range list {
			a.NotEqual(item, str)
		}
		list = append(list, str)
	}
}
