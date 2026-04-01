package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"ppt-probe/src/watcher"
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

	watcher.Execute(targetFile, *outputDir)
	fmt.Printf("Probing: %s\nSaving to: %s\n", targetFile, *outputDir)

	if *isWatch {
		watcher.WatchFile(context.Background(), targetFile, *outputDir)
	}

	fmt.Println("Done!")
}
