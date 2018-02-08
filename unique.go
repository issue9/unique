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

const formatBase = 10

var defaultUnique = New(time.Now().Unix(), 2, 60)

// Unique 基于时间戳的唯一字符串。
//
// NOTE: 不能在一秒之内重置计数器。
type Unique struct {
	random *rand.Rand

	prefix   string
	timer    *time.Timer
	duration int64

	ai   *autoinc.AutoInc
	step int64
}

// New 声明一个新的 Unique。
//
// seed 随机种子；
// step 计数器的最大步长，可以负数，为 0 会 panic；
// duration 计数器的最长重置时间，单位秒。系统会在 [1,timer] 范围内重置计数器；
func New(seed, step, duration int64) *Unique {
	random := rand.New(rand.NewSource(seed))

	u := &Unique{
		step:     step,
		duration: duration,
		random:   random,
	}

	u.reset()

	return u
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	u.prefix = strconv.FormatInt(time.Now().Unix(), formatBase)

	if u.ai != nil {
		u.ai.Stop()
	}
	u.ai = autoinc.New(1, u.getRandomNumber(u.step), 1000)

	resetTime := time.Duration(u.getRandomNumber(u.duration)) * time.Minute
	if u.timer != nil {
		u.timer.Stop()
	}
	u.timer = time.AfterFunc(resetTime, u.reset)
}

// String 返回一个唯一的字符串
func (u *Unique) String() string {
	id, err := u.ai.ID()
	if err == nil {
		return u.prefix + strconv.FormatInt(id, formatBase)
	}

	u.reset()

	return u.prefix + strconv.FormatInt(u.ai.MustID(), formatBase)
}

// Bytes 返回 String() 的 []byte 格式
func (u *Unique) Bytes() []byte {
	return []byte(u.String())
}

// 获取一个位于 [1,max) 区间的值
func (u *Unique) getRandomNumber(max int64) int64 {
	n := u.random.Int63n(max)
	if n <= 0 {
		n++
	}

	return n
}

// String 返回一个唯一的字符串
func String() string {
	return defaultUnique.String()
}

// Bytes 返回 String() 的 []byte 格式
func Bytes() []byte {
	return defaultUnique.Bytes()
}
