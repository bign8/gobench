package web

//go:generate appcfg.py update .

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/", index)
}

var (
	indexTPL = template.Must(template.ParseFiles("index.html"))
	static   = http.FileServer(http.Dir("static"))
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		static.ServeHTTP(w, r)
		return
	}
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	indexTPL.Execute(w, map[string]interface{}{
		"user": u.String(),
	})
}

func upload(w http.ResponseWriter, r *http.Request) {
	// TODO: process post
}
