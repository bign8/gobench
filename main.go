package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const version = "v0.0.0"

var (
	in  = flag.String("in", "", "input file (default to stdin)")
	out = flag.String("out", "", "output file (default to stdout)")
	ver = flag.Bool("version", false, "print version and exit")
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
	log.SetFlags(0) // No time ... prefix for error messages
	if *ver {
		fmt.Printf("gobench %s\n", version)
		os.Exit(0)
	}

	// Actually process output
	inp, outp, err := getIO()
	if err == nil {
		scanner := bufio.NewScanner(inp)
		for scanner.Scan() && err == nil {
			line := scanner.Text()
			if strings.HasPrefix(line, "Benchmark") {
				err = parseLine(line, outp)
			}
		}
		if err == nil {
			err = scanner.Err()
		}
	}
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}

func parseLine(line string, o io.Writer) (err error) {
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
	return json.NewEncoder(o).Encode(res)
}
