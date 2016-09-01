package main

import (
	"html/template"
	"net/http"
	"strings"

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
		"/plotly.js":   true,
		"/favicon.ico": true,
		"/humans.txt":  true,
		"/robots.txt":  true,
	}
)

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: fan this stuff out (datastore requests that is)
	if forward[r.URL.Path] {
		static.ServeHTTP(w, r) // NOTE: appengine automagically gzips content on supporting clients
		return
	}

	vars := make(map[string]interface{})
	loc := r.URL.Path[1:]
	vars["path"] = loc

	if strings.Contains(loc, "/Bench:") {

		// Fetch benchmark
		var ben bench
		parts := strings.Split(loc, "/Bench:")
		key := datastore.NewKey(ctx, "Bench", parts[1], 0, path2key(ctx, parts[0]))
		if err := datastore.Get(ctx, key, &ben); err != nil {
			if err == datastore.ErrNoSuchEntity {
				http.Redirect(w, r, "..", http.StatusSeeOther)
				return
			}
			log.Warningf(ctx, "Bench: %s", err)
		}
		vars["bench"] = ben

		// Fetch points
		q := datastore.NewQuery("Point").Ancestor(key).Order("stamp")
		var points []point
		if _, err := q.GetAll(ctx, &points); err != nil {
			log.Warningf(ctx, "Points: %s", err)
		}
		vars["points"] = toPoint(points)

	} else {
		parent := path2key(ctx, loc)

		// Fetch sub-paths
		q := datastore.NewQuery("Path").Filter("parent =", parent).Order("name")
		var paths []path
		if _, err := q.GetAll(ctx, &paths); err != nil {
			log.Warningf(ctx, "Path: %s", err)
		}
		vars["children"] = paths

		// Fetch benchmarks
		if parent != nil {
			q = datastore.NewQuery("Bench").Ancestor(parent).Order("name")
			var benches []bench
			if _, err := q.GetAll(ctx, &benches); err != nil {
				log.Warningf(ctx, "Benches: %s", err)
			}
			vars["benches"] = benches
		}
	}

	// User session handler
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

	// Navigation handler
	pat := strings.Split(loc, "/")
	nav := make([]breadcrumb, len(pat))
	var base string
	for i, sec := range pat {
		base = base + "/" + sec
		nav[i] = breadcrumb{
			Name: sec,
			Link: base,
		}
	}
	vars["nav"] = nav

	// TODO: handle the case if nothing is found
	// TODO: otherwise cache response (probably a wrapper function)
	// log.Errorf(ctx, "Page Not Found: %q %s", r.URL.Path, parent)
	// http.NotFound(w, r)
	// return
	if err := indexTPL.Execute(w, vars); err != nil {
		log.Errorf(ctx, "Executing Template: %s", err)
	}
}
