package main

//go:generate appcfg.py update .

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", route)
}

func main() {
	appengine.Main()
}

func route(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	start := time.Now()
	defer func() {
		if re := recover(); re != nil {
			log.Criticalf(ctx, "Problem: %#v %s", re, debug.Stack())
			http.Error(w, "FAIL WHALE!", http.StatusInternalServerError)
		} else {
			log.Debugf(ctx, "%s %s: %s", r.Method, r.URL, time.Since(start))
		}
	}()
	// TODO: OPTIONS
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	switch r.Method {
	case "GET":
		index(ctx, w, r)
	case "POST":
		upload(ctx, w, r)
	default:
		// Technically http.StatusMethodNotAllowed, but less visible errors
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func path2key(ctx context.Context, path string) *datastore.Key {
	key := rootKey(ctx)
	for _, folder := range strings.Split(path, "/") {
		key = datastore.NewKey(ctx, "Path", folder, 0, key)
	}
	return key
}

func key2path(key *datastore.Key) string {
	var parts []string
	for key != nil {
		parts = append(parts, key.StringID())
		key = key.Parent()
	}
	parts = parts[:len(parts)-1] // strip off root key
	// Reverse and join
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "/")
}

func rootKey(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Root", "root", 0, nil)
}
