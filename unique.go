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

var defaultUnique = New(time.Now().Unix(), 2, 60, "20060102150405-", false)

// Unique 基于时间戳的唯一字符串。
//
// Unique 由两部分组成：
// 前缀是由一个相对稳定的字符串，与时间相关联；
// 后缀是一个自增的数值。
// 每次刷新前缀之后，都会重置后缀的计数器，从头开始。
// 刷新时间和计数器的步长都是一个随机数。
//
// NOTE: 不能在一秒之内重置计数器。
type Unique struct {
	random     *rand.Rand
	formatBase int

	prefix       string
	prefixFormat string
	timer        *time.Timer
	duration     int64

	step int64
	ai   *autoinc.AutoInc
}

// New 声明一个新的 Unique。
//
// seed 随机种子；
// step 计数器的最大步长，只能大于 0；
// duration 计数器的最长重置时间，单位秒。系统会在 [1,duration] 范围内重置计数器；
// prefixFormat 格式化 prefix 的方式，若指定，则格式化为时间，否则将时间戳转换为数值。
// alpha 是否包含字母
func New(seed, step, duration int64, prefixFormat string, alpha bool) *Unique {
	random := rand.New(rand.NewSource(seed))

	if step <= 0 {
		panic("无效的参数 step")
	}

	if duration <= 0 {
		panic("无效的参数 duration")
	}

	u := &Unique{
		random:       random,
		formatBase:   10,
		duration:     duration,
		prefixFormat: prefixFormat,
		step:         step,
	}

	if alpha {
		u.formatBase = 36
	}

	u.reset()

	return u
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	if u.prefixFormat != "" {
		u.prefix = time.Now().Format(u.prefixFormat)
	} else {
		u.prefix = strconv.FormatInt(time.Now().Unix(), u.formatBase)
	}

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
		return u.prefix + strconv.FormatInt(id, u.formatBase)
	}

	u.reset()

	return u.prefix + strconv.FormatInt(u.ai.MustID(), u.formatBase)
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
