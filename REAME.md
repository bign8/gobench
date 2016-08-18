# GoBench
Converts go unit test files to usable and trackable format.

## Usage
```sh
go test -race -bench=. -benchmem -v ./... | tee test.out
gobench -in test.out
```
