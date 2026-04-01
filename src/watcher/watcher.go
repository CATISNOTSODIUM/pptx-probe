package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const POLL_TIME_SECONDS = 2

func WatchFile(ctx context.Context, filePath string, outputDir string) error {
	fmt.Println("Watch for changes")
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	ticker := time.NewTicker(POLL_TIME_SECONDS * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			stat, err := os.Stat(filePath)
			if err != nil {
				// We log instead of returning so a temporary
				// filesystem hiccup doesn't crash the watcher.
				log.Printf("Error stating file %s: %v", filePath, err)
				continue
			}

			if stat.Size() != initialStat.Size() || !stat.ModTime().Equal(initialStat.ModTime()) {
				fmt.Printf("Change detected in %s. Applying updates...\n", filePath)

				Execute(filePath, outputDir)

				// Update initialStat with the current stat to avoid another syscall
				initialStat = stat
			}
		}
	}
}
