package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// defaults for testing
// var regex = regexp.MustCompile("^(.*?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
// var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	var dryRun bool
	var regexInput, replaceString string
	flag.BoolVar(&dryRun, "dry", true, "whether or not this should be a dry run (i.e. won't actually rename)")
	flag.StringVar(&regexInput, "regex", "", "the regex to match the files you want to rename")
	flag.StringVar(&replaceString, "replaceString", "", "the replacement filenames you want to output")
	flag.Parse()
	regex := regexp.MustCompile(regexInput)

	dir := "sample"
	var toRename []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name(), regex, replaceString); err == nil {
			toRename = append(toRename, path)
		}
		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFilename, _ := match(filename, regex, replaceString)
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

func match(filename string, regex *regexp.Regexp, replaceString string) (string, error) {
	if !regex.MatchString(filename) {
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}
	return regex.ReplaceAllString(filename, replaceString), nil
}
