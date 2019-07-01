package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type file struct {
	name string
	path string
}

func main() {
	dir := "sample"
	var toRename []file

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, file{
				name: info.Name(),
				path: path,
			})
		}
		return nil
	})
	for _, f := range toRename {
		fmt.Printf("%q\n", f)
	}
	for _, orig := range toRename {
		var n file
		var err error
		n.name, err = match(orig.name)
		if err != nil {
			fmt.Println("Error matching: ", orig.path, err.Error())
		}
		n.path = filepath.Join(dir, n.name)
		fmt.Printf("mv %s => %s\n", orig.path, n.path)

		err = os.Rename(orig.path, n.path)
		if err != nil {
			fmt.Println("Error renaming: ", orig.path, err.Error())
		}
	}
	// files, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	panic(err)
	// }
	// count := 0
	// var toRename []string

	// for _, file := range files {
	// 	if file.IsDir() {
	// 		fmt.Println("Dir: ", file.Name())
	// 	} else {
	// 		_, err := match(file.Name(), 0)
	// 		if err == nil {
	// 			count++
	// 			toRename = append(toRename, file.Name())
	// 		}
	// 	}
	// }
}

// match returns the new filename or an error if filename didn't match our pattern
func match(filename string /*, total int*/) (string, error) {
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
	//return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
	return fmt.Sprintf("%s - %d.%s", strings.Title(name), number, ext), nil
}
