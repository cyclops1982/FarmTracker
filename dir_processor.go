package main

import (
	"log"
	"flag"
	"path/filepath"
	"os"
	"time"
)

func FindFiles(ch chan string, path string) {
	foundDirs := make(map[string]bool)
	for {
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if info.IsDir() == false {
				if foundDirs[p] == false {
					ch <- p
					foundDirs[p] = true
				}
			}
			return nil
		})
		time.Sleep(1 * time.Second)
	}
}


var inputdir string
func main() {
	flag.StringVar(&inputdir, "inputdir", "dumps/", "The directory to process.")
	flag.Parse()

	ch := make(chan string)
	go FindFiles(ch, inputdir)

	for file := range ch {
		log.Println("Processing ", file)
	}

}
