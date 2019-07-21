// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unique

import "github.com/issue9/rands"

// Rands 生成唯一的随机字符串
//
// 相对于 Unique，Rands 更有随机性。
type Rands struct {
	// 固定的前缀内容
	//
	// 如果不需要，可以为空
	Prefix string

	Unique *Unique
	Rands  *rands.Rands
}

// Bytes 返回 []byte 内容
func (r *Rands) Bytes() []byte {
	return []byte(r.String())
}

// String 返回字符串内容
func (r *Rands) String() string {
	return r.Prefix + r.Unique.String() + r.Rands.String()
}
