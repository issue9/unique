// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unique

import "github.com/issue9/rands"

// Rands 生成唯一的随机字符串
//
// 相对于 Unqiue，Rands 更有随机性。
type Rands struct {
	Unique *Unique
	Rands  *rands.Rands
}

// Bytes 返回 []byte 内容
func (r *Rands) Bytes() []byte {
	return append(r.Unique.Bytes(), r.Rands.Bytes()...)
}

// String 返回字符串内容
func (r *Rands) String() string {
	return r.Unique.String() + r.Rands.String()
}
