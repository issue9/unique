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

// Unique 基于时间戳的唯一字符串。
// 不能在一秒之内重置计数器。
type Unique struct {
	step   int64 // 计数器的最大步长
	timer  int64 // 计数器的最大重置时间
	random *rand.Rand

	prefix string           // 时间戳字符串
	ai     *autoinc.AutoInc // 自增值
}

// New 声明一个新的 Unique。
// seed 随机种子；
// step 计数器的最大步长；
// timer 计数器的最长重置时间，单位秒
func New(seed, step, timer int64) *Unique {
	random := rand.New(rand.NewSource(seed))

	u := &Unique{
		step:   step,
		timer:  timer,
		random: random,
		ai:     autoinc.New(1, random.Int63n(step), 1000),
	}

	u.reset()

	return u
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	u.prefix = strconv.FormatInt(time.Now().Unix(), 10)
	u.ai.Reset(1, u.random.Int63n(u.step))

	restTime := time.Duration(u.random.Int63n(u.timer)) * time.Minute
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

var defaultUnique = New(time.Now().Unix(), 10, 60)

// String 返回一个唯一的 ID
func String() string {
	return defaultUnique.String()
}

// Bytes 返回 String() 的 []byte 格式
func Bytes() []byte {
	return defaultUnique.Bytes()
}
