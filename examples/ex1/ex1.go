package main

import (
	"fmt"
	"regexp"

	"github.com/dustywilson/dir/directory"
	"github.com/dustywilson/dir/file"
)

func main() {
	root := directory.CreateDirectory(nil, "My Root Directory")

	directory.CreateDirectory(root, "Fruit").CreateDirectory("Green").CreateDirectory("Grapes").CreateDirectory("Sour")
	file.CreateFile(directory.CreateDirectory(root, "Meat").CreateDirectory("Cow"), "Ground Chuck")
	cake := directory.CreateDirectory(root, "Dessert").CreateDirectory("Cake")

	file.CreateFile(cake, "Angel Food")

	dirs, err := root.FindDirectories(regexp.MustCompile(`[Rr]`), -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nMatched %d directories.\n", len(dirs))
	for _, d := range dirs {
		fmt.Println(d)
	}

	files, err := root.FindFiles(regexp.MustCompile(`[Oo]`), -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nMatched %d files.\n", len(files))
	for _, f := range files {
		fmt.Println(f)
	}
}
