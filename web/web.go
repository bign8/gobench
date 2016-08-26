package main

//go:generate appcfg.py update .

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/cloud/datastore"
)

func init() {
	http.HandleFunc("/", route)
}

func main() {
	appengine.Main()
}

func route(w http.ResponseWriter, r *http.Request) {
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
