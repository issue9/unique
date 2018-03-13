// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package unique 产生一个唯一字符串
package unique

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/issue9/autoinc"
)

var stringInst, numberInst, dateInst *Unique

// Unique 基于时间戳的唯一字符串，长度不固定。
//
// NOTE: 算法是基于系统时间的。所以必须得保证时间上正确的，否则可能会造成非唯一的情况。
//
// Unique 由两部分组成：
// 前缀是由一个相对稳定的字符串，与时间相关联；
// 后缀是一个自增的数值。
//
// 每次刷新前缀之后，都会重置后缀的计数器，从头开始。
// 刷新时间和计数器的步长都是一个随机数。
type Unique struct {
	random *rand.Rand

	// 数据转换成字符串所采用的进制。
	formatBase int

	// 前缀部分的内容。
	//
	// 根据 prefixFormat 是否存在，会呈现不同的内容：
	// 如果 prefixFormat 为空，prefix 为一个时间戳的整数值，
	// 按一定的进制进行转换之后的值；否则是按 prefixFormat
	// 进行格式化的时间数据。
	prefix       string
	prefixFormat string

	timer    *time.Timer
	duration int64

	step int64
	ai   *autoinc.AutoInc

	// 用保证 prefix 和 ai 的一致性。
	resetLocker sync.RWMutex
}

// String 初始化一个以字符串形式表示唯一值的 Unique 实例，大致格式如下：
//  p4k5f81
//
// NOTE: 多次调用，返回的是同一个实例。
func String() *Unique {
	if stringInst == nil {
		stringInst = New(time.Now().Unix(), 2, 60, "", 36)
	}

	return stringInst
}

// Number 初始化一个数字形式表示唯一值的 Unique 实例，大致格式如下：
//  15193130121
//
// NOTE: 多次调用，返回的是同一个实例。
func Number() *Unique {
	if numberInst == nil {
		numberInst = New(time.Now().Unix(), 2, 60, "", 10)
	}

	return numberInst
}

// Date 初始化一个以日期形式表示唯一值的 Unique 实例，大致格式如下：
//  20180222232332-1
//
// NOTE: 多次调用，返回的是同一个实例。
func Date() *Unique {
	if dateInst == nil {
		dateInst = New(time.Now().Unix(), 2, 60, "20060102150405-", 10)
	}

	return dateInst
}

// New 声明一个新的 Unique。
//
// seed 随机种子；
// step 计数器的最大步长，只能大于 0；
// duration 计数器的最长重置时间，单位秒。系统会在 [1,duration] 范围内重置计数器；
// prefixFormat 格式化 prefix 的方式，若指定，则格式化为时间，否则将时间戳转换为数值。
// base 数值转换成字符串时，所采用的进制，可以是 [2,36] 之间的值。
func New(seed, step, duration int64, prefixFormat string, base int) *Unique {
	if step <= 0 {
		panic("无效的参数 step")
	}

	if duration <= 0 {
		panic("无效的参数 duration")
	}

	if prefixFormat != "" && !isValidDateFormat(prefixFormat) {
		panic("无效的 prefixFormat 参数")
	}

	if base < 2 || base > 36 {
		panic("无效的 base 值，只能介于 [2,36] 之间")
	}

	u := &Unique{
		random:       rand.New(rand.NewSource(seed)),
		formatBase:   base,
		duration:     duration,
		prefixFormat: prefixFormat,
		step:         step,
	}

	u.reset()

	return u
}

func isValidDateFormat(format string) bool {
	return strings.Contains(format, "2006") &&
		strings.Contains(format, "01") &&
		strings.Contains(format, "02") &&
		strings.Contains(format, "15") &&
		strings.Contains(format, "04") &&
		strings.Contains(format, "05")
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	u.resetLocker.Lock()
	defer u.resetLocker.Unlock()

	if u.prefixFormat != "" {
		u.prefix = time.Now().Format(u.prefixFormat)
	} else {
		u.prefix = strconv.FormatInt(time.Now().Unix(), u.formatBase)
	}

	if u.ai != nil {
		go u.ai.Stop()
	}
	u.ai = autoinc.New(1, u.getRandomNumber(u.step), 1000)

	if u.timer != nil {
		u.timer.Stop()
	}
	dur := time.Duration(u.getRandomNumber(u.duration)) * time.Minute
	u.timer = time.AfterFunc(dur, u.reset)
}

// String 返回一个唯一的字符串
func (u *Unique) String() string {
	u.resetLocker.RLock()
	p := u.prefix
	id, ok := u.ai.ID()
	u.resetLocker.RUnlock()

	for !ok {
		u.reset() // NOTE: reset 包含对 resetLocker 的操作

		u.resetLocker.RLock()
		p = u.prefix
		id, ok = u.ai.ID()
		u.resetLocker.RUnlock()
	}

	return p + strconv.FormatInt(id, u.formatBase)
}

// Bytes 返回 String() 的 []byte 格式
//
// 在多次出错之后，可能会触发 panic
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
