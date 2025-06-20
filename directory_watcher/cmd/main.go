package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Directory Monitor")

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatalf(`Error creating watcher: %v`, err)
		return
	}
	defer watcher.Close()

	notifications := make(chan string)

	go func() {
		for notification := range notifications {
			fmt.Println(notification)
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				timestamp := time.Now().Format(time.RFC3339)

				if event.Op.Has(fsnotify.Create) {
					notifications <- timestamp + " File created:" + event.Name
				}

				if event.Op.Has(fsnotify.Write) {
					notifications <- timestamp + " File updated:" + event.Name
				}

				if event.Op.Has(fsnotify.Remove) {
					notifications <- timestamp + " File removed:" + event.Name
				}

				if event.Op.Has(fsnotify.Rename) {
					notifications <- timestamp + " File renamed:" + event.Name
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatal("Error:", err)
			}
		}
	}()

	err = watcher.Add("../watched_directory")
	if err != nil {
		log.Fatalf("Error adding directory to watcher: %v", err)
	}

	fmt.Println("Watching directory: ./watched_directory")

	// Block forever
	<-make(chan struct{})
}
