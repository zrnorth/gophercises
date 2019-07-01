package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry", true, "whether or not this should be a dry run (i.e. won't actually rename)")
	flag.Parse()

	dir := "sample"
	toRename := make(map[string][]string) // storing a map of filepaths to lists of files

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		curDir := filepath.Dir(path)
		if m, err := match(info.Name()); err == nil {
			key := filepath.Join(curDir, fmt.Sprintf("%s.%s", m.base, m.ext))
			toRename[key] = append(toRename[key], info.Name())
		}
		return nil
	})

	// Get the total number of files so the "m of n" string will have the correct n
	for key, files := range toRename {
		dir := filepath.Dir(key)
		n := len(files)
		sort.Strings(files)
		for i, filename := range files {
			res, _ := match(filename)
			newFilename := fmt.Sprintf("%s - %d of %d.%s", res.base, (i + 1), n, res.ext)
			oldPath := filepath.Join(dir, filename)
			newPath := filepath.Join(dir, newFilename)
			fmt.Printf("mv %s => %s\n", oldPath, newPath)

			// Actually do the rename
			if !dryRun {
				err := os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Println("Error renaming: ", oldPath, err.Error())
				}
			}
		}
	}
}

type matchResult struct {
	base  string
	index int
	ext   string
}

func match(filename string) (*matchResult, error) {
	// split into "birthday", "001", "txt"
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	fileNameWithoutPeriods := strings.Join(pieces[0:len(pieces)-1], ".")

	pieces = strings.Split(fileNameWithoutPeriods, "_")
	name := strings.Join(pieces[0:len(pieces)-1], ".")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s didn't match our pattern", filename)
	}
	return &matchResult{strings.Title(name), number, ext}, nil
}
