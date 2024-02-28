// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package unique

import (
	"context"

	"github.com/issue9/rands/v2"
)

// Rands 生成唯一的随机字符串
//
// 相较于 [Unique] 此对象更具有随机性。
type Rands struct {
	// 固定的前缀内容
	//
	// 如果不需要，可以为空
	Prefix string

	Unique *Unique
	Rands  *rands.Rands
}

func (r *Rands) Serve(ctx context.Context) error {
	go r.Unique.Serve(ctx)
	go r.Rands.Serve(ctx)

	<-ctx.Done()
	return ctx.Err()
}

// Bytes 返回 []byte 内容
func (r *Rands) Bytes() []byte { return []byte(r.String()) }

// String 返回字符串内容
func (r *Rands) String() string {
	return r.Prefix + r.Unique.String() + r.Rands.String()
}
