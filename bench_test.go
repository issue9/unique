// SPDX-FileCopyrightText: 2024 caixw
//
// SPDX-License-Identifier: MIT

package unique

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/issue9/rands/v3"
)

func BenchmarkUnique(b *testing.B) {
	s := NewString(10)
	ctx, cancel := context.WithCancel(context.Background())
	go s.Serve(ctx)
	time.Sleep(time.Microsecond * 500)
	defer cancel()

	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}

func BenchmarkRands(b *testing.B) {
	r := rand.New(rand.NewPCG(0, 0))
	s := NewRands(10, r, 10, 11, rands.AlphaNumber())

	ctx, cancel := context.WithCancel(context.Background())
	go s.Serve(ctx)
	time.Sleep(time.Microsecond * 500)
	defer cancel()

	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}
