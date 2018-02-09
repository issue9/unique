// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unique

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() {
		New(20, 0, 100, "", false)
	})

	a.Panic(func() {
		New(20, 1, 0, "", false)
	})

	a.Panic(func() {
		New(20, 1, 1, "2006", false)
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

	u := New(time.Now().Unix(), 2, 5, "", true)
	a.NotNil(u)

	// 保证 getRandomNumber 不会返回 0
	for i := 0; i <= 100; i++ {
		a.Equal(u.getRandomNumber(1), 1)
	}
}
