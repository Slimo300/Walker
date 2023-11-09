package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Walker struct {
	LinesTotal       int
	ExtensionCounter map[string]int

	AcceptedExtensions    map[string]bool
	AcceptedExtensionsSet bool

	OmitBlank bool
}

func (w *Walker) WithAcceptedExtensions(exts []string) {
	extensionsMap := make(map[string]bool)

	for _, ext := range exts {
		extensionsMap[fmt.Sprintf(".%s", ext)] = true
	}

	w.AcceptedExtensions = extensionsMap
	w.AcceptedExtensionsSet = true
}

func (w *Walker) WithOmitBlank() {
	w.OmitBlank = true
}

func (res *Walker) Print() {
	log.Println("Total number of lines in directory: ", res.LinesTotal)

	for ext, num := range res.ExtensionCounter {
		log.Printf("%s%d\n", tabulate(ext), num)
	}
}

func (w *Walker) CountLines(dir string) error {

	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if w.AcceptedExtensionsSet && !w.AcceptedExtensions[filepath.Ext(path)] {
			return nil
		}

		fileReader, err := os.Open(path)
		if err != nil {
			return err
		}

		linesCount := 0

		scanner := bufio.NewScanner(fileReader)
		for scanner.Scan() {
			if w.OmitBlank && len(strings.TrimSpace(scanner.Text())) == 0 {
				continue
			}
			linesCount += 1
		}

		w.LinesTotal += linesCount

		if _, ok := w.ExtensionCounter[filepath.Ext(path)]; !ok {
			w.ExtensionCounter[filepath.Ext(path)] = linesCount
		} else {
			w.ExtensionCounter[filepath.Ext(path)] += linesCount
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func main() {

	acceptedExtensions := flag.String("ext", "", "comma separated extensions of files you want to count ")
	omitBlank := flag.Bool("omit-blank", false, "boolean value stating whether to count in blank lines")

	flag.Parse()

	dir := flag.Args()[0]
	if dir == "" {
		log.Fatal("Path cannot be blank")
	}

	walker := Walker{
		ExtensionCounter: map[string]int{},
	}

	if len(*acceptedExtensions) > 0 {
		walker.WithAcceptedExtensions(strings.Split(*acceptedExtensions, ","))
	}
	if *omitBlank {
		walker.WithOmitBlank()
	}

	if err := walker.CountLines(dir); err != nil {
		log.Fatal(err)
	}

	walker.Print()
}

func tabulate(str string) string {
	if len(str) > 3 {
		return fmt.Sprintf("%s\t", str)
	} else {
		return fmt.Sprintf("%s\t\t", str)
	}
}
