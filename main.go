package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO: post metadata about the build

const (
	version = "v0.0.0"

	// FAIL	_/home/miki/Projects/goroot/src/xunit	0.004s
	// ok  	_/home/miki/Projects/goroot/src/anotherTest	0.000s
	gtSuiteRE = "^(?:ok|FAIL)[ \t]+([^ \t]+)[ \t]+" //(\\d+.\\d+)"
)

var (
	in  = flag.String("in", "", "input file (default to stdin)")
	out = flag.String("out", "", "output file (default to stdout)")
	ver = flag.Bool("version", false, "print version and exit")
	url = flag.String("host", os.Getenv("GOBENCH_URL"), "url of the gobench server")

	findSuite = regexp.MustCompile(gtSuiteRE).FindStringSubmatch
)

func getIn() (io.Reader, error) {
	if *in == "" {
		return os.Stdin, nil
	}
	return os.Open(*in)
}

func getOut() (io.Writer, error) {
	if *out == "" {
		return os.Stdout, nil
	}
	return os.Create(*out)
}

func getIO() (inp io.Reader, outp io.Writer, err error) {
	inp, err = getIn()
	if err == nil {
		outp, err = getOut()
	}
	return inp, outp, err
}

func main() {
	flag.Parse()
	if *ver {
		fmt.Fprintf(os.Stdout, "gobench %s\n", version)
		os.Exit(0)
	}
	inp, outp, err := getIO()

	// re-arrange outp if we have a url to post to
	var buf bytes.Buffer
	if *url != "" {
		outp = io.MultiWriter(outp, &buf)
	}

	// Process Output
	parser := batcher{Encoder: json.NewEncoder(outp)}
	if err == nil {
		scanner := bufio.NewScanner(inp)
		for scanner.Scan() && err == nil {
			err = parser.parse(scanner.Text())
		}
		if err == nil {
			err = scanner.Err()
		}
	}

	// Post to gobench server if possible
	if err == nil && *url != "" {
		var resp *http.Response
		resp, err = http.Post(*url+"/upload", "application/json", &buf)
		if err == nil {
			fmt.Fprintf(os.Stderr, "Response (%s): ", resp.Status)
			_, err = io.Copy(os.Stderr, resp.Body)
			resp.Body.Close()
		}
		if err == nil && resp.StatusCode != http.StatusOK {
			err = errors.New("failed to store data")
		}
	}

	// Kill if errors are not good
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

type batcher struct {
	*json.Encoder
	list []map[string]interface{}
}

func (b *batcher) parse(line string) (err error) {
	if tokens := findSuite(line); tokens != nil {
		err = b.flush(tokens[1])
	} else if strings.HasPrefix(line, "Benchmark") {
		err = b.parseLine(line)
	}
	return err
}

func (b *batcher) flush(suite string) error {
	for _, obj := range b.list {
		obj["suite"] = suite
		if err := b.Encode(obj); err != nil {
			return err
		}
	}
	b.list = nil // reset list
	return nil
}

func (b *batcher) parseLine(line string) (err error) {
	cols := strings.Split(line[9:], "\t") // 9 = len(benchmark)
	for i, c := range cols {
		cols[i] = strings.TrimSpace(c)
	}
	res := map[string]interface{}{
		"name": strings.SplitN(cols[0], "-", 2)[0], // Clip CPU from Benchmark Name
	}
	res["iter"], err = strconv.Atoi(cols[1])
	if err == nil {
		for _, k := range cols[2:] {
			parts := strings.SplitN(k, " ", 2)
			res[parts[1]], err = strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return err
			}
		}
	}
	b.list = append(b.list, res)
	return nil
}
