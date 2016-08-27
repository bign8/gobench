package main

import (
	"html/template"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

var (
	indexTPL = template.Must(template.ParseFiles("index.html"))
	static   = http.FileServer(http.Dir("static"))
	forward  = map[string]bool{
		"/css.css":     true,
		"/js.js":       true,
		"/favicon.ico": true,
	}
)

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if forward[r.URL.Path] {
		static.ServeHTTP(w, r)
		return
	} else if r.URL.Path != "/" {
		parent := path2key(ctx, r.URL.Path[1:])
		// TODO: parse full dirty path
		log.Errorf(ctx, "Page Not Found: %q %s", r.URL.Path, parent)
		http.NotFound(w, r)
		return
	}

	// Home page handler
	u := user.Current(ctx)
	vars := make(map[string]interface{})
	if u != nil {
		vars["user"] = u.String()
	}
	indexTPL.Execute(w, vars)
}
