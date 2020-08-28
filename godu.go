package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"syscall"

	"github.com/karrick/godirwalk"
)

func out(fullsize int64, fullname string) {
	fsize := float64(fullsize) / 1024.0 / 1024.0
	result := int64(math.Ceil(fsize))
	fmt.Printf("%d\t%s\n", result, fullname)
}

func main() {
	dirname := "."
	if len(os.Args) > 1 {
		dirname = os.Args[1]
	}

	fullsize, err := du(dirname)
	if err != nil {
		log.Printf("Failed to du \"%s\": %s", dirname, err)
		os.Exit(1)
	}
	out(fullsize, dirname)
}

func du(dirname string) (int64, error) {
	var fullsize int64
	scanner, err := godirwalk.NewScanner(dirname)
	if err != nil {
		return fullsize, fmt.Errorf("Failed to create scanner for directory \"%s\": %s", dirname, err)
	}
	for scanner.Scan() {
		dirent, err := scanner.Dirent()
		if err != nil {
			log.Printf("Failed to get directory entry: %s", err)
			continue
		}
		fullname := fmt.Sprintf("%s%c%s", dirname, os.PathSeparator, dirent.Name())
		if dirent.IsDir() && !dirent.IsSymlink() {
			dirsize, err := du(fullname)
			if err != nil {
				return fullsize, err
			}
			out(dirsize, fullname)
			fullsize += dirsize
			continue
		}
		var info syscall.Stat_t
		err = syscall.Lstat(fullname, &info)
		if err != nil {
			return fullsize, fmt.Errorf("Failed to lstat(\"%s\"): %s", fullname, os.PathError{"lstat", fullname, err})
		}
		filesize := info.Blocks * 512
		fullsize += filesize
		if false {
			out(filesize, fullname)
		}
	}
	return fullsize, nil
}
