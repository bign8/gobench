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

	w.Header().Add("X-Clacks-Overhead", "GNU Terry Pratchett")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")

	switch r.Method {
	case http.MethodGet:
		index(ctx, w, r)
	case http.MethodPost:
		upload(ctx, w, r)
	case http.MethodOptions: // NO-OP
	case http.MethodHead: // NO-OP
	default:
		http.Error(w, "Method \""+r.Method+"\" not allowed", http.StatusMethodNotAllowed)
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
