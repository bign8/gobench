package main

import (
	"html/template"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
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
	}

	// TODO: fan this stuff out
	vars := make(map[string]interface{})
	parent := path2key(ctx, r.URL.Path[1:])

	// Fetch sub-paths
	q := datastore.NewQuery("Path").Filter("parent =", parent).Order("name")
	var paths []path
	_, err := q.GetAll(ctx, &paths)
	log.Infof(ctx, "shiz: %#v %#v", parent, paths)
	if err != nil {
		log.Errorf(ctx, "Error w/Path query: %s", err)
	}
	vars["children"] = paths

	// Home page handler
	u := user.Current(ctx)
	if u != nil {
		out, _ := user.LogoutURL(ctx, "/")
		vars["user"] = map[string]string{
			"name":   u.String(),
			"logout": out,
		}
	} else {
		vars["login"], _ = user.LoginURL(ctx, "/")
	}

	// TODO: handle the case if nothing is found
	// log.Errorf(ctx, "Page Not Found: %q %s", r.URL.Path, parent)
	// http.NotFound(w, r)
	// return
	if err := indexTPL.Execute(w, vars); err != nil {
		log.Errorf(ctx, "Executing Template: %s", err)
	}
}
