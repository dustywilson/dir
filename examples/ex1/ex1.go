package main

import (
	"fmt"
	"regexp"

	"github.com/dustywilson/dir"
	"github.com/dustywilson/dir/directory"
)

func main() {
	root := directory.CreateDirectory(nil, "My Root Directory")

	directory.CreateDirectory(root, "Fruit").CreateDirectory("Green").CreateDirectory("Grapes").CreateDirectory("Sour")
	directory.CreateDirectory(root, "Grains").CreateDirectory("Wheat")
	directory.CreateDirectory(root, "Dessert").CreateDirectory("Cake")

	dirs, err := root.FindDirectories(regexp.MustCompile(`[Rr]`), -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Matched %d directories.\n", len(dirs))
	for _, d := range dirs {
		fmt.Println(dir.Path(d))
	}
}
