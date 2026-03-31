package main

import (
	"flag"
	"fmt"
	"os"
	"ppt-probe/src/extractor"
	"ppt-probe/src/models"
)

func main() {
	outputDir := flag.String("o", "output", "Directory to save extracted code")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	args := flag.Args()
	if *help || len(args) < 1 {
		fmt.Println("Usage: ./pptx-probe [options] <path-to-presentation.pptx>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	targetFile := args[0]
	fmt.Printf("Probing: %s\nSaving to: %s\n", targetFile, *outputDir)

	ppt, err := models.ReadPowerPoint(targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not read PPTX: %v\n", err)
		os.Exit(1)
	}

	for i, bytes := range ppt.Slides {
		xmlNode, err := models.Decode(bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to decode slide %s: %v\n", i, err)
			continue
		}

		extractor.Parse(*xmlNode, *outputDir)
	}

	fmt.Println("Done!")
}
