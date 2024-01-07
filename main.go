/*
Package main provides a command line tool for filtering JSON lines.

// https://github.com/tidwall/gjson/blob/master/SYNTAX.md
It uses GJSON query syntax, which is designed to query JSON arrays;
this program will ephemerally convert each line is not already a JSON array
in order to be able to use this syntax.
If the query result for the line (optionally wrapped into an array) is Existing,
the original line will be printed to stdout.
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

var flagMatchAll = flag.String("match-all", "", "match all of these properties (gjson syntax, comma separated queries)")
var flagMatchAny = flag.String("match-any", "", "match any of these properties (gjson syntax, comma separated queries)")
var flagMatchNone = flag.String("match-none", "", "match none of these properties (gjson syntax, comma separated queries)")
var errInvalidMatchAll = errors.New("invalid match-all")
var errInvalidMatchAny = errors.New("invalid match-any")
var errInvalidMatchNone = errors.New("invalid match-none")
var errInvalidLine = errors.New("invalid line (is it valid JSON?)")

func filterStream(reader io.Reader, writer io.Writer, matchAll []string, matchAny []string, matchNone []string) {
	breader := bufio.NewReader(reader)
	bwriter := bufio.NewWriter(writer)

readLoop:
	for {
		read, err := breader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, os.ErrClosed) || errors.Is(err, io.EOF) {
				break
			}
			log.Fatalln(err)
		}
		if err := filter(read, matchAll, matchAny, matchNone); err != nil {
			if errors.Is(err, errInvalidLine) {
				log.Fatalln(err)
			}
			continue readLoop
		}
		bwriter.Write(read)
		bwriter.Flush()
	}
}

// filter filters some read line on the matchAll, matchAny, and matchNone queries.
// These queries should be written in GJSON query syntax.
// https://github.com/tidwall/gjson/blob/master/SYNTAX.md
func filter(read []byte, matchAll []string, matchAny []string, matchNone []string) error {

	parsed := gjson.ParseBytes(read)
	if !parsed.Exists() {
		return errInvalidLine
	}

	// Here we hack the line into an array containing only this datapoint.
	// This allows us to use the GJSON query syntax, which is designed for use with arrays, not single objects.
	if !parsed.IsArray() {
		read = []byte(fmt.Sprintf("[%s]", string(read)))
	}

	for _, query := range matchAll {
		if res := gjson.GetBytes(read, query); !res.Exists() {
			return fmt.Errorf("%w: %s", errInvalidMatchAll, query)
		}
	}

	didMatchAny := len(matchAny) == 0
	for _, query := range matchAny {
		if gjson.GetBytes(read, query).Exists() {
			didMatchAny = true
			break
		}
	}
	if !didMatchAny {
		return fmt.Errorf("%w: %s", errInvalidMatchAny, matchAny)
	}

	for _, query := range matchNone {
		if gjson.GetBytes(read, query).Exists() {
			return fmt.Errorf("%w: %s", errInvalidMatchNone, query)
		}
	}
	return nil
}

func splitFlagStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func main() {
	filterStream(os.Stdin, os.Stdout, splitFlagStringSlice(*flagMatchAll), splitFlagStringSlice(*flagMatchAny), splitFlagStringSlice(*flagMatchNone))
}
