// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package unique 产生一个唯一字符串
package unique

import (
	"context"
	"strconv"
	"strings"
	"time"
)

// Unique 基于时间戳的唯一不定长字符串
//
// Unique 由两部分组成：
// 前缀是由一个相对稳定的字符串，与时间相关联；
// 后缀是一个自增的数值。
//
// 每次刷新前缀之后，都会重置后缀的计数器，从头开始。
// 刷新时间和计数器的步长都是一个随机数。
//
// NOTE: 算法是基于系统时间的。所以必须得保证时间上正确的，否则可能会造成非唯一的情况。
// NOTE: 产生的数据有一定的顺序规则。
type Unique struct {
	// 数值转换成字符串所采用的进制
	//
	// 同时应用于自增的数值后缀，以及 prefix 的非时间格式。
	formatBase int

	// 前缀部分的内容
	//
	// 根据 prefixFormat 是否存在，会呈现不同的内容：
	//  - 空值，prefix 为时间戳，按 formatBase 进制进行转换之后的字符串；
	//  - 非空，按 prefixFormat 进行格式化的时间格式。
	prefix       string
	prefixFormat string

	channel chan string

	timer    *time.Timer
	duration time.Duration
}

// NewString 声明以字符串形式表示的 [Unique] 实例
//
// 格式为：p4k5f81
func NewString(bufferSize int) *Unique { return New(bufferSize, time.Hour, "", 36) }

// NewNumber 声明以数字形式表示的 [Unique] 实例
//
// 格式为：15193130121
func NewNumber(bufferSize int) *Unique { return New(bufferSize, time.Hour, "", 10) }

// NewDate 声明以日期形式表示的 [Unique] 实例
//
// 格式为：20180222232332-1
func NewDate(bufferSize int) *Unique {
	return New(bufferSize, time.Hour, "20060102150405-", 10)
}

// New 声明一个新的 [Unique]
//
// 每一秒，最多能产生 [math.MaxInt64] 个唯一值，需求量超过此值的不适合。
//
// bufferSize 缓存大小，不能小于 1；
// duration 计数器的重置时间，不能小于 1*time.Second；
// prefixFormat 格式化 prefix 的方式，若指定，则格式化为时间，否则将时间戳转换为数值；
// base 数值转换成字符串时，所采用的进制，可以是 [2,36] 之间的值。
func New(bufferSize int, duration time.Duration, prefixFormat string, base int) *Unique {
	if bufferSize < 1 {
		panic("参数 bufferSize 不能小于 1")
	}
	if duration < time.Second {
		panic("参数 duration 不能小于 1 秒")
	}

	if prefixFormat != "" && !isValidDateFormat(prefixFormat) {
		panic("无效的 prefixFormat 参数")
	}

	if base < 2 || base > 36 {
		panic("参数 base 只能介于 [2,36] 之间")
	}

	return &Unique{
		formatBase:   base,
		duration:     duration,
		prefixFormat: prefixFormat,

		channel: make(chan string, bufferSize),
	}
}

func isValidDateFormat(format string) bool {
	return strings.Contains(format, "2006") &&
		strings.Contains(format, "01") &&
		strings.Contains(format, "02") &&
		strings.Contains(format, "15") &&
		strings.Contains(format, "04") &&
		strings.Contains(format, "05")
}

func (u *Unique) Serve(ctx context.Context) error {
	u.reset(ctx)

	<-ctx.Done()
	u.timer.Stop()
	return ctx.Err()
}

// 重置时间戳和计数器
func (u *Unique) reset(ctx context.Context) {
	if u.prefixFormat != "" {
		u.prefix = time.Now().Format(u.prefixFormat)
	} else {
		u.prefix = strconv.FormatInt(time.Now().Unix(), u.formatBase)
	}

	cc, cancel := context.WithCancel(ctx)
	go func() {
		var n int64
		for {
			select {
			case <-cc.Done():
				return
			default:
				n++
				u.channel <- u.prefix + strconv.FormatInt(n, u.formatBase)
			}
		}
	}()

	if u.timer != nil {
		u.timer.Stop()
		u.timer = nil
	}
	u.timer = time.AfterFunc(u.duration, func() {
		cancel()
		u.reset(ctx)
	})
}

// String 返回一个唯一的字符串
func (u *Unique) String() string { return <-u.channel }

// Bytes 返回 [Unique.String] 的 []byte 格式
//
// 在多次出错之后，可能会触发 panic
func (u *Unique) Bytes() []byte { return []byte(u.String()) }
