// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package unique 产生一个唯一字符串
package unique

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/issue9/autoinc"
)

const (
	resetTimeMax = 60 // 分种
	resetStepMax = 10
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

// Unique 基于时间戳的唯一字符串。
// 不能在一秒之内重置计数器。
type Unique struct {
	prefix string           // 时间戳字符串
	ai     *autoinc.AutoInc // 自增值
}

// New 声明一个新的 Unique
func New() *Unique {
	step := int64(random.Intn(resetStepMax))

	u := &Unique{
		ai: autoinc.New(1, step, 1000),
	}

	u.reset()

	return u
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	u.prefix = strconv.FormatInt(time.Now().Unix(), 10)
	u.ai.Reset(1, int64(random.Intn(resetStepMax)))

	restTime := time.Duration(random.Intn(resetTimeMax)) * time.Minute
	time.AfterFunc(restTime, u.reset)
}

// String 返回一个唯一的 ID
func (u *Unique) String() string {
	return u.prefix + strconv.FormatInt(u.ai.MustID(), 10)
}

// Bytes 返回 String() 的 []byte 格式
func (u *Unique) Bytes() []byte {
	return []byte(u.String())
}

var defaultUnique = New()

// String 返回一个唯一的 ID
func String() string {
	return defaultUnique.String()
}

// Bytes 返回 String() 的 []byte 格式
func Bytes() []byte {
	return defaultUnique.Bytes()
}
