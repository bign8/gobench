package main

import (
	"time"

	"google.golang.org/cloud/datastore"
)

type jsonBench struct {
	Suite string  `json:"suite"`
	Name  string  `json:"name"`
	N     uint64  `json:"iter"`
	NS    float64 `json:"ns/op"`
	B     uint64  `json:"B/op"`
	Alloc uint64  `json:"allocs/op"`
}

// DB Key: Type: Path, String: Name, Ancestor: Path (parent of current directory)
type path struct {
	Name string `datastore:"name"`
}

// DB Key: Type: Bench, String: Name, Ancestor: Path (full parent path as ancestors)
type bench struct {
	Name string    `datastore:"name"`
	Seen time.Time `datastore:"seen"`
}

// DB Key: Type: Point, String/Int: Generated, Ancestor: Bench (with full ancestor tree)
type point struct {
	N     uint64         `datastore:"iter"`
	NS    float64        `datastore:"ns"`
	B     uint64         `datastore:"b"`
	Alloc uint64         `datastore:"allocs"`
	Batch *datastore.Key `datastore:"batch"`
}

// DB Key: Type: Batch, String/Int: Generated, Ancestor: Path (top most directory of the benchmarks)
type batch struct {
	Stamp    time.Time `datastore:"stamp"`
	Branch   string    `datastore:"branch"`
	Commit   string    `datastore:"commit"`
	Build    uint      `datastore:"build"`
	BuildURL string    `datastore:"build_url"`
	Tag      string    `datastore:"tag"`
	Slug     string    `datastore:"slug"`
	// Codecov other: yaml, service, flags, pr, job
}

type trie map[string]*trie

func (t *trie) Add(parts []string) {
	if len(parts) == 0 {
		return
	}
	child, ok := (*t)[parts[0]]
	if !ok {
		child = &trie{}
		(*t)[parts[0]] = child
	}
	child.Add(parts[1:])
}

func (t *trie) Prefix() []string {
	if len(*t) != 1 {
		return nil
	}
	var key string
	for key = range *t {
		/* nothing */
	}
	return append([]string{key}, (*t)[key].Prefix()...)
}
