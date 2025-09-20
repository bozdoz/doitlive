package main

import (
	_ "embed"
	"time"

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
			debug("added", path)
		}
		return nil
	})

	var name string

	// dedupe
	t := time.AfterFunc(24*time.Hour, func() {
		fmt.Println("[doitlive] Change:", name)
		broadcast <- "reload"
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			debug(event)
			if event.Has(fsnotify.Write) {
				// dedupe
				name = event.Name
				t.Reset(200 * time.Millisecond)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("[doitlive] Watcher error:", err)
		}
	}
}
