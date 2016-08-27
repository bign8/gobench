package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func upload(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parse entire payload
	var err error
	var set []*point
	dec := json.NewDecoder(r.Body)
	for err == nil {
		var ben point
		if err = dec.Decode(&ben); err == nil {
			set = append(set, &ben)
		}
	}
	if err != io.EOF || len(set) == 0 {
		http.Error(w, "Invalid Data Format", http.StatusExpectationFailed)
		return
	}

	// TODO (bign8): place batch metadata and get batch key
	var batchKey *datastore.Key

	// Determine set of benches needed to be stored
	now := time.Now()
	tree := make(trie)
	benches := make([]*pair, 0, len(set))
	for _, ben := range set {
		tree.Add(strings.Split(ben.Suite, "/"))
		ben.Batch = batchKey
		benches = append(benches, &pair{
			Key: datastore.NewKey(ctx, "Bench", ben.Name, 0, path2key(ctx, ben.Suite)),
			Val: &bench{Name: ben.Name, Seen: now},
		})
	}
	getAll(ctx, benches) // Don't care about errs (just pre-fetching so we don't lose attributes later)
	newBenches := make([]*pair, 0, len(set))
	for _, p := range benches {
		if p.Val.(*bench).Seen == now {
			newBenches = append(newBenches, p)
		}
	}

	// walk the tree and build all the path objects
	objs := tree.Walk(&pair{
		Key: rootKey(ctx),
	}, func(val string, parent interface{}) interface{} {
		return &pair{
			Key: datastore.NewKey(ctx, "Path", val, 0, parent.(*pair).Key),
			Val: &path{Name: val},
		}
	})
	paths := make([]*pair, len(objs))
	for i, obj := range objs {
		paths[i] = obj.(*pair)
	}

	// iterate paths and generate the points
	points := make([]*pair, 0, len(set))
	for _, pat := range set {
		parent := datastore.NewKey(ctx, "Bench", pat.Name, 0, path2key(ctx, pat.Suite))
		points = append(points, &pair{
			Key: datastore.NewIncompleteKey(ctx, "Point", parent),
			Val: pat,
		})
	}
	paths = append(paths, points...)

	// Lets store all this data and return to the user
	err = putAll(ctx, append(paths, newBenches...))
	if err != nil {
		log.Criticalf(ctx, "Error Storing Data: %s", err)
		http.Error(w, "Problem Storing Data", http.StatusInternalServerError)
	} else {
		slug := strings.Join(tree.Prefix(3), "/") // <host>/<user>/<repo>
		log.Infof(ctx, "Full Suite Prefix: %q", slug)
		fmt.Fprintf(w, "Success! Avilable at http://%s/%s\n", r.Host, slug)
	}
}

func putAll(ctx context.Context, list []*pair) error {
	keys := make([]*datastore.Key, len(list))
	vals := make([]interface{}, len(list))
	for i, p := range list {
		keys[i] = p.Key
		vals[i] = p.Val
	}
	_, err := datastore.PutMulti(ctx, keys, vals)
	return err
}

func getAll(ctx context.Context, list []*pair) error {
	keys := make([]*datastore.Key, len(list))
	vals := make([]interface{}, len(list))
	for i, p := range list {
		keys[i] = p.Key
		vals[i] = p.Val
	}
	return datastore.GetMulti(ctx, keys, vals)
}
