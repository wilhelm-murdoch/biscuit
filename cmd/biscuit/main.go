package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/docopt/docopt.go"
	"github.com/wilhelm-murdoch/biscuit"
)

var (
	version = "Biscuit 0.0.1"
	usage   = `Biscuit - A simple command line utility for language detection.

	Usage:
	 biscuit (--file=<file> | --text=<text>) [--length=<length>] [--extension=<extension>] [--table] --path=<path> 
	 biscuit (--help | --version)

	Examples:
	 biscuit -t "c'est la vie." -p /path/to/corpora
	 biscuit -f /path/to/test/file -p /path/to/corpora
	 biscuit --file /path/to/test/file --path /path/to/corpora --extension .data

	Options:
	 -f --file=<file>            Path to file containing text you want to match.
	 -t --text=<text>            The text you want to match.
	 -l --length=<length>        The desired Ngram length. [default: 3]
	 -p --path=<path>            Path to library of comparison texts.
	 -e --extension=<extension>  File extension for comparison texts. [default: .txt]
	 --table                     Will display a table of sorted results.
	 -h --help                   Will display this help screen.
	 -v --version                Displays the current version of Biscuit.`
)

func main() {
	arguments, err := docopt.Parse(usage, nil, true, version, false)

	if err != nil {
		fmt.Println("Could not properly execute command; exiting ...")
		os.Exit(1)
	}

	content := arguments["--text"].(string)
	if len(content) == 0 {
		file := arguments["--file"]
		if file != nil {
			bytes, err := ioutil.ReadFile(file.(string))

			if err != nil {
				fmt.Println("Invalid file specified:", err)
				os.Exit(1)
			}

			content = string(bytes)
		}
	}

	if len(content) == 0 {
		fmt.Println("There is nothing to score ...")
		os.Exit(1)
	}

	path := filepath.Clean(filepath.Dir(arguments["--path"].(string)))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("There was a problem with your path:", err)
		os.Exit(1)
	}

	pathInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("There was a problem with your path:", err)
		os.Exit(1)
	}

	if !pathInfo.IsDir() {
		fmt.Println("--path must be a directory ...")
		os.Exit(1)
	}

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("There was a problem with your path:", err)
		os.Exit(1)
	}

	ext := arguments["--extension"].(string)
	var files []string
	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == strings.ToLower(ext) {
			files = append(files, path+"/"+file.Name())
		}
	}

	if len(files) == 0 {
		fmt.Println("Could not find comparison library. Double check your path and extension.")
		os.Exit(1)
	}

	l := arguments["--length"].(string)
	length, err := strconv.Atoi(l)
	if err != nil {
		fmt.Println("There was a problem defining ngram length:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	profiles := make([]*biscuit.Profile, 0, len(files))

	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			fileName := filepath.Base(file)
			label := fileName[0 : len(fileName)-len(filepath.Ext(fileName))]

			profile, err := biscuit.NewProfileFromFile(label, file, length)

			if err == nil {
				profiles = append(profiles, profile)
			}

		}(file)
	}

	wg.Wait()

	unknown := biscuit.NewProfileFromText("unknown", content, length)

	table := arguments["--table"].(bool)
	if table {
		sorted, scores, _ := unknown.MatchReturnAll(profiles)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 9, 5, ' ', 0)

		fmt.Fprintln(w, "Label:\tScore:")

		for _, key := range sorted {
			score := biscuit.Round(100*scores[key], 3)
			fmt.Fprintln(w, fmt.Sprintf("%s\t", key)+fmt.Sprintf("%f", score)[0:4]+"%")
		}

		w.Flush()
	} else {
		match, _ := unknown.MatchReturnBest(profiles)
		fmt.Println(match)
	}

	os.Exit(0)
}
