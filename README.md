# GoBench
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/bign8/gobench.svg)](https://travis-ci.org/bign8/gobench)
[![codecov](https://codecov.io/gh/bign8/gobench/branch/master/graph/badge.svg)](https://codecov.io/gh/bign8/gobench)
[![Go Report Card](https://goreportcard.com/badge/github.com/bign8/gobench)](https://goreportcard.com/report/github.com/bign8/gobench)
[![GoDoc](https://godoc.org/github.com/bign8/gobench?status.svg)](https://godoc.org/github.com/bign8/gobench)
[![GitHub release](http://img.shields.io/github/release/bign8/pipelines.svg)](https://github.com/bign8/pipelines/releases)

Converts go unit test files to usable and trackable format.

## Usage
```sh
go test -bench=. -benchmem -run=# ./... | gobench
# or
go test -bench=. -benchmem -run=# ./... | tee test.out
gobench -in test.out
```

### Notes
* The `-race` parameter slows down benchmarks substantially.
* If tests have been run via a previous command, they can be skipped using the `-run` flag with an invalid name token such at `-run=#`.
  ```
  go test -bench=. -benchmem -run=# ./...
  ```
