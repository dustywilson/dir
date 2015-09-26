package main

import (
	"fmt"
	"regexp"

	"github.com/dustywilson/directory/directory"
)

func main() {
	root, err := directory.CreateDirectory(nil, "My Root Directory")
	if err != nil {
		panic(err)
	}

	directory.CreateDirectory(root, "Fruit")
	directory.CreateDirectory(root, "Grains")
	directory.CreateDirectory(root, "Dessert")

	dirs, err := root.FindDirectories(regexp.MustCompile(`[Rr]`), -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Matched %d directories.\n", len(dirs))
	for _, dir := range dirs {
		fmt.Printf("%+v\n", dir)
	}
}
