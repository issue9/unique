// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unique

import (
	"math"
	"testing"
	"time"

	"github.com/issue9/autoinc"

	"github.com/issue9/assert"
)

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
		New(20, 1, 1, "20060102150405", 2)
	})
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

func TestUnique_getRandomNumber(t *testing.T) {
	a := assert.New(t)

	u := New(time.Now().Unix(), 2, 5, "", 36)
	a.NotNil(u)

	// 保证 getRandomNumber 不会返回 0
	for i := 0; i <= 100; i++ {
		a.Equal(u.getRandomNumber(1), 1)
	}
}

func TestUnique_String(t *testing.T) {
	a := assert.New(t)

	u := New(time.Now().Unix(), 100000, 50, "", 5)
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

	u := New(time.Now().Unix(), 100000, 50, "", 5)
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
