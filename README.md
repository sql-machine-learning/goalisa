# goalisa: A Go `database/sql` driver for alisa

[![Build Status](https://travis-ci.org/sql-machine-learning/goalisa.svg?branch=develop)](https://travis-ci.org/sql-machine-learning/goalisa) 
[![GoDoc](https://godoc.org/github.com/sql-machine-learning/goalisa?status.svg)](https://godoc.org/github.com/sql-machine-learning/goalisa) 
[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](LICENSE) 
[![Go Report Card](https://goreportcard.com/badge/github.com/sql-machine-learning/goalisa)](https://goreportcard.com/report/github.com/sql-machine-learning/goalisa)

# What is goalisa
To access databases, Go programmers call the standard library `database/sql`, which relies on drivers to talk to database management systems. `goalisa` is such a driver that talks to alisa.

# For Users
`goalisa` is go-gettable. Please run the following command to install it:

```bash
go get sqlflow.org/goalisa
```
`sqlflow.org/goalisa` is a [vainty import path](https://blog.bramp.net/post/2017/10/02/vanity-go-import-paths/) of `goalisa`.

Please make sure you have Go 1.13+.

# License
`goalisa` comes with [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).
