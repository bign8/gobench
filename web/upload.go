package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	ctx := appengine.NewContext(r)

	// Parse entire payload
	var err error
	var set []*jsonBench
	for err == nil {
		var ben jsonBench
		if err = dec.Decode(&ben); err == nil {
			log.Debugf(ctx, "Data: %#v", ben)
			set = append(set, &ben)
		}
	}

	// TODO: put data in datastore

	if err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintln(w, "Success! (TODO: Return Post location)")
	}
}
