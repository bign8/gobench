package main

//go:generate appcfg.py update .

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/", index)
}

func main() {
	appengine.Main()
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
	vars := make(map[string]interface{})
	if u != nil {
		vars["user"] = u.String()
	}
	indexTPL.Execute(w, vars)
}

type bench struct {
	Suite string  `json:"suite"`
	Name  string  `json:"name"`
	N     uint64  `json:"iter"`
	NS    float64 `json:"ns/op"`
	B     uint64  `json:"B/op"`
	Alloc uint64  `json:"allocs/op"`
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	ctx := appengine.NewContext(r)

	// Parse entire payload
	var err error
	var set []*bench
	for err == nil {
		var ben bench
		if err = dec.Decode(&ben); err == nil {
			log.Debugf(ctx, "Data: %#v", ben)
			set = append(set, &ben)
		}
	}

	if err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintln(w, "Success! (TODO: Return Post location)")
	}
}
