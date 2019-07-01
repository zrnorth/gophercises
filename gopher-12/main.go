package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := "birthday_001.txt"
	// convert into "Birthday - 1 of 4.txt"
	newName, err := match(filename)
	if err != nil {
		fmt.Println("no match")
		os.Exit(0)
	}
	fmt.Println(newName)
}

// match returns the new filename or an error if filename didn't match our pattern
func match(filename string) (string, error) {
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
	return fmt.Sprintf("%s - %d.%s", strings.Title(name), number, ext), nil
}
