package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var source = `?   	github.com/bign8/games	[no test files]
?   	github.com/bign8/games/cmd/chess	[no test files]
?   	github.com/bign8/games/cmd/game-server	[no test files]
?   	github.com/bign8/games/cmd/games	[no test files]
=== RUN   TestLocationParseGood
--- PASS: TestLocationParseGood (0.00s)
=== RUN   TestRookMoves
--- PASS: TestRookMoves (0.00s)
=== RUN   TestStateIsCheck
--- PASS: TestStateIsCheck (0.00s)
=== RUN   TestParse
--- PASS: TestParse (0.00s)
=== RUN   TestFENFullCycle
--- PASS: TestFENFullCycle (0.00s)
=== RUN   TestNewString
--- PASS: TestNewString (0.00s)
PASS
BenchmarkStateApply-8    	 2000000	       864 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocationParse-8 	10000000	       174 ns/op	       0 B/op	       0 allocs/op
BenchmarkLocationOffset-8	100000000	        19.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMoves-8         	   30000	     56528 ns/op	    2720 B/op	       2 allocs/op
BenchmarkMovesClip-8     	   30000	     48760 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesRook-8     	 1000000	      1171 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesPawn-8     	10000000	       117 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesKnight-8   	 2000000	       699 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesBishop-8   	 2000000	       737 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesQueen-8    	 1000000	      1858 ns/op	       0 B/op	       0 allocs/op
BenchmarkMovesKing-8     	 2000000	       874 ns/op	      64 B/op	       2 allocs/op
BenchmarkStateIsCheck-8  	 1000000	      1900 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseFEN-8      	 1000000	      1063 ns/op	     224 B/op	       2 allocs/op
BenchmarkStringFEN-8     	 1000000	      2281 ns/op	      66 B/op	       3 allocs/op
BenchmarkStateString-8   	   50000	     34273 ns/op	    5568 B/op	      27 allocs/op
ok  	github.com/bign8/games/impl/chess	30.451s
PASS
BenchmarkStateString-8	  500000	      2583 ns/op	     240 B/op	       1 allocs/op
BenchmarkNewState-8   	 1000000	      1936 ns/op	     360 B/op	       6 allocs/op
BenchmarkApply-8      	 3000000	       511 ns/op	     128 B/op	       2 allocs/op
BenchmarkTerminal-8   	20000000	        65.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkUtility-8    	20000000	        81.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkSVG-8        	  200000	      8013 ns/op	    5056 B/op	       6 allocs/op
ok  	github.com/bign8/games/impl/ttt	11.140s
?   	github.com/bign8/games/player/cli	[no test files]
?   	github.com/bign8/games/player/minimax	[no test files]`

var expected = `{"B/op":128,"allocs/op":1,"iter":2000000,"name":"StateApply","ns/op":864,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":10000000,"name":"LocationParse","ns/op":174,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":100000000,"name":"LocationOffset","ns/op":19.2,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":2720,"allocs/op":2,"iter":30000,"name":"Moves","ns/op":56528,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":30000,"name":"MovesClip","ns/op":48760,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":1000000,"name":"MovesRook","ns/op":1171,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":10000000,"name":"MovesPawn","ns/op":117,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":2000000,"name":"MovesKnight","ns/op":699,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":2000000,"name":"MovesBishop","ns/op":737,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":1000000,"name":"MovesQueen","ns/op":1858,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":64,"allocs/op":2,"iter":2000000,"name":"MovesKing","ns/op":874,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":0,"allocs/op":0,"iter":1000000,"name":"StateIsCheck","ns/op":1900,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":224,"allocs/op":2,"iter":1000000,"name":"ParseFEN","ns/op":1063,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":66,"allocs/op":3,"iter":1000000,"name":"StringFEN","ns/op":2281,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":5568,"allocs/op":27,"iter":50000,"name":"StateString","ns/op":34273,"suite":"github.com/bign8/games/impl/chess"}
{"B/op":240,"allocs/op":1,"iter":500000,"name":"StateString","ns/op":2583,"suite":"github.com/bign8/games/impl/ttt"}
{"B/op":360,"allocs/op":6,"iter":1000000,"name":"NewState","ns/op":1936,"suite":"github.com/bign8/games/impl/ttt"}
{"B/op":128,"allocs/op":2,"iter":3000000,"name":"Apply","ns/op":511,"suite":"github.com/bign8/games/impl/ttt"}
{"B/op":0,"allocs/op":0,"iter":20000000,"name":"Terminal","ns/op":65.4,"suite":"github.com/bign8/games/impl/ttt"}
{"B/op":0,"allocs/op":0,"iter":20000000,"name":"Utility","ns/op":81.3,"suite":"github.com/bign8/games/impl/ttt"}
{"B/op":5056,"allocs/op":6,"iter":200000,"name":"SVG","ns/op":8013,"suite":"github.com/bign8/games/impl/ttt"}`

func TestMain(m *testing.M) {
	flag.Parse()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HIJACKED!")
	}))
	*url = s.URL
	out := m.Run()
	s.Close()
	os.Exit(out)
}

func TestRun(t *testing.T) {
	var out bytes.Buffer
	in := bytes.NewBufferString(source)
	err := run(in, &out)
	if err != nil {
		t.Error(err)
	}

	got := json.NewDecoder(bytes.NewReader(out.Bytes()))
	has := json.NewDecoder(bytes.NewReader([]byte(expected)))

	var er1, er2 error
	var src, dst map[string]interface{}
	for ; er1 == nil && er2 == nil; er1, er2 = got.Decode(&dst), has.Decode(&src) {
		if !reflect.DeepEqual(src, dst) {
			t.Fatalf("%v != %v", src, dst)
		}
	}
	if er1 != er2 {
		t.Error(er1)
		t.Error(er2)
	}
}

func TestGetIO(t *testing.T) {
	*in, *out = "", ""
	i, o, e := getIO()
	if e != nil {
		t.Fatal(e)
	}
	if i != os.Stdin {
		t.Fatal("Input Doesn't equal os.Stdin")
	}
	if o != os.Stdout {
		t.Fatal("Ouput Doesn't equal os.Stdout")
	}
	open = func(string) (*os.File, error) { return nil, nil }
	create = func(string) (*os.File, error) { return nil, nil }
	*in, *out = "IN", "OUT"
	_, _, e = getIO()
	if e != nil {
		t.Fatal(e)
	}
}

func BenchmarkParseLine(b *testing.B) {
	var o batcher
	for i := 0; i < b.N; i++ {
		o.parseLine("BenchmarkLocationOffset-8	100000000	        19.2 ns/op	       0 B/op	       0 allocs/op")
	}
}
