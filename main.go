package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"ppt-probe/src/watcher"
	"syscall"
)

func main() {
	outputDir := flag.String("o", "output", "Directory to save extracted code")
	help := flag.Bool("h", false, "Show help")
	isWatch := flag.Bool("w", false, "watch for changes")

	flag.Parse()

	args := flag.Args()
	if *help || len(args) < 1 {
		fmt.Println("Usage: ./pptx-probe [options] <path-to-presentation.pptx>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	targetFile := args[0]

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	watcher.Execute(targetFile, *outputDir)
	log.Printf("Probing: %s\nSaving to: %s\n", targetFile, *outputDir)

	if *isWatch {
		err := watcher.WatchFile(ctx, targetFile, *outputDir)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Done!")
}
