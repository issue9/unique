// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package unique

import (
	"context"
	"math/rand/v2"
	"unsafe"

	"github.com/issue9/rands/v3"
)

// Rands 生成唯一的随机字符串
//
// [Rands] 由两部分组成：
//   - [Unique] 负责保证字符串的唯一性，但是内容是有序的；
//   - [rands.Rands] 负责生成混淆的随机字符；
//
// 两者结合可以保证生成的内容唯一且无序。相较于 [Unique] 此对象更具有随机性。
type Rands struct {
	unique *Unique
	rands  *rands.Rands[byte]
}

// NewRands 声明 [Rands]
//
// bufferSize 负责初始化 [Unique] 和 [rands.Rands] 对象；
// r, min, max, bs 负责初始化 [rands.Rands] 对象；
func NewRands(bufferSize int, r *rand.Rand, min, max int, bs []byte) *Rands {
	return &Rands{
		unique: NewString(bufferSize),
		rands:  rands.New(r, bufferSize, min, max, bs),
	}
}

func (r *Rands) Serve(ctx context.Context) error {
	go r.unique.Serve(ctx)
	go r.rands.Serve(ctx)

	<-ctx.Done()
	return ctx.Err()
}

// Bytes 返回 []byte 内容
func (r *Rands) Bytes() []byte {
	s := r.String()
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// String 返回字符串内容
func (r *Rands) String() string { return r.unique.String() + r.rands.String() }
