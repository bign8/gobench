package main

//go:generate appcfg.py update .

import (
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/datastore"
)

func init() {
	http.HandleFunc("/", route)
}

func main() {
	appengine.Main()
}

func route(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			log.Criticalf(appengine.NewContext(r), "Problem: %#v %s", re, debug.Stack())
			http.Error(w, "FAIL WHALE!", http.StatusInternalServerError)
		}
	}()
	// TODO: defer recover here
	// TODO: OPTIONS
	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	switch r.Method {
	case "GET":
		index(w, r)
	case "POST":
		upload(w, r)
	default:
		// Technically http.StatusMethodNotAllowed, but less visible errors
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func path2key(ctx context.Context, path string) (key *datastore.Key) {
	for _, folder := range strings.Split(path, "/") {
		key = datastore.NewKey(ctx, "Path", folder, 0, key)
	}
	return key
}
