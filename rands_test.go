// SPDX-FileCopyrightText: 2024-2026 caixw
//
// SPDX-License-Identifier: MIT

package unique

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/issue9/assert/v5"
	"github.com/issue9/rands/v3"
)

func TestRands(t *testing.T) {
	a := assert.New(t, false)
	r := rand.New(rand.NewPCG(0, 0))

	u := NewRands(10, r, 10, 20, rands.AlphaNumber())
	ctx, cancel := context.WithCancel(context.Background())
	go u.Serve(ctx)
	time.Sleep(time.Microsecond * 500)
	defer cancel()

	list := make([]string, 0, 100)
	for range 100 {
		time.Sleep(time.Millisecond * 20)
		str := u.String()
		for _, item := range list {
			a.NotEqual(item, str)
		}
		list = append(list, str)
	}
}
