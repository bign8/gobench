# GoBench
[![Build Status](https://travis-ci.org/bign8/gobench.svg)](https://travis-ci.org/bign8/gobench)

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
