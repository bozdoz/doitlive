package main

import (
	_ "embed"

	"fmt"
	"io/fs"
	"path/filepath"

	// file watcher
	"github.com/fsnotify/fsnotify"
)

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer watcher.Close()

	// Watch every subdirectory in root
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if _, found := ignored[path]; found {
				return filepath.SkipDir
			}
			if err := watcher.Add(path); err != nil {
				fmt.Printf("%v", err)
			}
		}
		return nil
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("[doitlive] Change:", event.Name)
				broadcast <- "reload"
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("[doitlive] Watcher error:", err)
		}
	}
}
