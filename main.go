package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

func getFileSize(fileInfo os.FileInfo) string {
	fileSize := fileInfo.Size()
	size := ""
	if fileSize > 0 {
		size = fmt.Sprintf("(%vb)", int(fileInfo.Size()))
	} else {
		size = "(empty)"
	}
	return size
}

func getVerticalPrefix(isLastItem bool) string {
	verticalLinePrefix := ""
	if isLastItem {
		verticalLinePrefix = ""
	} else {
		verticalLinePrefix = "│"
	}
	return verticalLinePrefix
}

func getPrefix(isLastItem bool) string {
	prefix := ""
	if isLastItem {
		prefix = "└───"
	} else {
		prefix = "├───"
	}
	return prefix
}

func renderDir(dir os.FileInfo, isLastItem bool, path string, parentPrefix string, printFiles bool) (string, error) {
	text := ""
	prefix := getPrefix(isLastItem)
	text += fmt.Sprintf( "%s%s%s\n", parentPrefix, prefix, dir.Name())
	newPath := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), dir.Name())
	verticalLinePrefix := getVerticalPrefix(isLastItem)
	childText, err := renderTree(newPath, printFiles, parentPrefix + verticalLinePrefix + "\t")
	if err != nil {
		return "", err
	}
	text += childText
	return text, nil
}

func renderFile(file os.FileInfo, prefix string) (string, error) {
	size := getFileSize(file)

	return fmt.Sprintf( "%s%s %s\n", prefix, file.Name(), size), nil
}

func renderTree(path string, printFiles bool, parentPrefix string) (string, error) {
	text := ""
	dirs, err := ioutil.ReadDir(path)

	if !printFiles {
		var items []os.FileInfo
		for _, item := range dirs {
			if item.IsDir() {
				items = append(items, item)
			}
		}
		dirs = make([]os.FileInfo, len(items), cap(items))
		copy(dirs, items)
	}

	sort.SliceStable(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name()	})
	if err != nil {
		return "", err
	}

	for index, item := range dirs {
		isLastItem := len(dirs) == index + 1
		prefix := getPrefix(isLastItem)

		if item.IsDir() {
			dirRow, renderDirErr := renderDir(item, isLastItem, path, parentPrefix, printFiles)
			if renderDirErr != nil {
				return "", renderDirErr
			}

			text += dirRow
		} else if printFiles {
			fileRow, renderFileErr := renderFile(item, parentPrefix + prefix)
			if renderFileErr != nil {
				return "", renderFileErr
			}

			text += fileRow
		}
	}

	return text, nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	text, err := renderTree(path, printFiles, "")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s", text)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
