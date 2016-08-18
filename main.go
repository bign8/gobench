package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Benchmark") {
				fmt.Fprintln(outp, line)
			}
		}
		err = scanner.Err()
	}
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
