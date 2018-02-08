// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unique

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestUnique_getRandomNumber(t *testing.T) {
	a := assert.New(t)

	u := New(time.Now().Unix(), 2, 5, true)
	a.NotNil(u)

	// 保证 getRandomNumber 不会返回 0
	for i := 0; i <= 100; i++ {
		a.Equal(u.getRandomNumber(1), 1)
	}
}
