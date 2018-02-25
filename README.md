unique [![Build Status](https://travis-ci.org/issue9/unique.svg?branch=master)](https://travis-ci.org/issue9/unique)
======


用于生成一个唯一字符串

```go
// 生成由数字和字母组成的唯一字符串，比如 p4k5f81
unique.String().String()

// 生成由数字组成的唯一字符串，比如 15193130121
unique.Number().String()

// 生成由日期与数字组成的唯一字符串，比如 20180222232332-1
unique.Date().String

// 或者可以自定义一个 Unique 实例
u := unique.New(time.Now().Unix(), 2, 60, "20060102150405-", 10)
u.String() // 生成唯一字符串。
```



### 安装

```shell
go get github.com/issue9/unique
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/unique)
[![GoDoc](https://godoc.org/github.com/issue9/unique?status.svg)](https://godoc.org/github.com/issue9/unique)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
