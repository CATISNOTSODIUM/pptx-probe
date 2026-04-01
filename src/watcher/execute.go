package watcher

import (
	"fmt"
	"os"
	"ppt-probe/src/extractor"
	"ppt-probe/src/models"
)

func Execute(filePath string, outputDir string) {
	ppt, err := models.ReadPowerPoint(filePath)
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

		extractor.Parse(*xmlNode, outputDir)
	}
}
