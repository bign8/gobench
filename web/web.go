package main

//go:generate appcfg.py update .

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/mjibson/appstats"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func init() {
	http.Handle("/", appstats.NewHandler(route))
}

func main() {
	appengine.Main()
}

func route(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func path2key(ctx context.Context, path string) (key *datastore.Key) {
	if path != "" {
		key = datastore.NewKey(ctx, "Path", path, 0, nil)
	}
	return key
}

func key2path(key *datastore.Key) string {
	return key.StringID()
}
