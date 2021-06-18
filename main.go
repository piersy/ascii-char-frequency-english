package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/piersy/gzip"
	"golang.org/x/net/html"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// This main will read the sgm files calculate the frequencies of ascii chars
// and save them to a a file.
func main() {
	err := aux()
	if err != nil {
		println(err.Error())
	}
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

}

func processFile(srcFile string, num int) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

}

func aux() error {
	countMap := make(map[byte]uint32)
	var count float64
	f, err := os.Open("reuters21578.tar.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	// gzf, err := gzip.NewReaderMaxBuf(f, 0xFFF)
	gzf, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to open gzip reader for file %s: %v", f.Name(), err)
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, ".sgm") {

			trans := transform.Chain(runes.Remove(runes.Predicate(func(r rune) bool { return r == 0x3 })))
			// io.Copy(os.Stdout, tarReader)
			t := html.NewTokenizer(tarReader)
			for {
				tt := t.Next()
				if tt == html.ErrorToken {
					break
				}

				name, _ := t.TagName()
				if tt == html.StartTagToken && string(name) == "body" {
					t.Next()
					s, _, _ := transform.String(trans, string(t.Text()))
					s = strings.TrimSpace(s)
					// Trim the Reuter final text, for some reason all the
					// articles seem to end with this, either uppercase or
					// capitalized.
					s = strings.TrimSuffix(s, "Reuter")
					s = strings.TrimSuffix(s, "REUTER")
					s = html.UnescapeString(s)
					for _, r := range s {
						if r > unicode.MaxASCII {
							fmt.Fprintf(os.Stderr, "Non ascii char found in %s %q   %q \n", header.Name, r, "\n")
							continue
						}
						countMap[byte(r)]++
						count++
					}
				}
			}
		}
	}

	type Result struct {
		Char byte
		Freq float64
	}
	results := make([]Result, 0, len(countMap))
	for k, v := range countMap {
		results = append(results, Result{byte(k), float64(v) / count})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Freq > results[j].Freq
	})
	println()
	println("Results, sorted highest freq first:")
	for _, r := range results {
		fmt.Printf("%q:%g\n", r.Char, r.Freq)
	}
	f, err = os.Create("ascii_freq.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	for _, r := range results {
		fmt.Fprintf(f, "%d:%g\n", r.Char, r.Freq)
	}
	bytes, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}
	println()
	println("Machine readable results written to ascii_freq.txt and ascii_freq.json")
	return ioutil.WriteFile("ascii_freq.json", bytes, 0666)

}
