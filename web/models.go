package main

type jsonBench struct {
	Suite string  `json:"suite"`
	Name  string  `json:"name"`
	N     uint64  `json:"iter"`
	NS    float64 `json:"ns/op"`
	B     uint64  `json:"B/op"`
	Alloc uint64  `json:"allocs/op"`
}
