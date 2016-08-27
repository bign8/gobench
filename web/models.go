package main

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type pair struct {
	Key *datastore.Key
	Val interface{}
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
	Suite string         `json:"suite" datastore:"-"`
	Name  string         `json:"name" datastore:"-"`
	N     int64          `json:"iter" datastore:"iter"`
	NS    float64        `json:"ns/op" datastore:"ns"`
	B     int64          `json:"B/op" datastore:"b"`
	Alloc int64          `json:"allocs/op" datastore:"allocs"`
	Batch *datastore.Key `json:"-" datastore:"batch"`
}

// DB Key: Type: Batch, String/Int: Generated, Ancestor: Path (top most directory of the benchmarks)
type batch struct {
	Stamp    time.Time `datastore:"stamp"`
	Branch   string    `datastore:"branch"`
	Commit   string    `datastore:"commit"`
	Build    uint      `datastore:"build"`
	BuildURL string    `datastore:"build_url"`
	Tag      string    `datastore:"tag"`
	// Codecov other: yaml, service, flags, pr, job, slug
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

func (t *trie) Prefix(max int) []string {
	if len(*t) != 1 || max == 0 {
		return nil
	}
	var key string
	for key = range *t {
		/* nothing */
	}
	return append([]string{key}, (*t)[key].Prefix(max-1)...)
}

type walker func(val string, parent interface{}) interface{}

func (t *trie) Walk(start interface{}, walk walker) []interface{} {
	var res []interface{}
	for key, value := range *t {
		parent := walk(key, start)
		res = append(res, parent)
		res = append(res, value.Walk(parent, walk)...)
	}
	return res
}
