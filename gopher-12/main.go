package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// filename := "birthday_001.txt"
	// // convert into "Birthday - 1 of 4.txt"
	// newName, err := match(filename)
	// if err != nil {
	// 	fmt.Println("no match")
	// 	os.Exit(0)
	// }
	// fmt.Println(newName)
	dir := "./sample"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	count := 0
	var toRename []string

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("Dir: ", file.Name())
		} else {
			_, err := match(file.Name(), 0)
			if err == nil {
				count++
				toRename = append(toRename, file.Name())
			}
		}
	}
	for _, origFilename := range toRename {
		origPath := filepath.Join(dir, origFilename)
		newFilename, err := match(origFilename, count)
		if err != nil {
			panic(err)
		}
		newPath := filepath.Join(dir, newFilename)
		fmt.Printf("mv %s => %s\n", origPath, newPath)
	}
	// origPath := fmt.Sprintf("%s/%s")
	// newPath := fmt.Sprintf("%s/%s")
}

// match returns the new filename or an error if filename didn't match our pattern
func match(filename string, total int) (string, error) {
	// split into "birthday", "001", "txt"
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	fileNameWithoutPeriods := strings.Join(pieces[0:len(pieces)-1], ".")

	pieces = strings.Split(fileNameWithoutPeriods, "_")
	name := strings.Join(pieces[0:len(pieces)-1], ".")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
