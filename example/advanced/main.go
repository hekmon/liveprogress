package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveprogress"
)

var (
	size3G = 3 * 1024 * 1024 * 1024
	size7G = 5 * 1024 * 1024 * 1024
	size8G = 8 * 1024 * 1024 * 1024
)

var (
	workers sync.WaitGroup
	spinner liveprogress.Spinner
)

func main() {
	// Global config (these are already the default values)
	liveprogress.Output = os.Stdout
	liveprogress.RefreshInterval = 100 * time.Millisecond
	// Progress bar configs
	liveprogress.DefaultConfig.Width = 40 // leave it a 0 for automatic width
	arrowsBarConfig := liveprogress.BarConfig{
		Width: 40, // leave it a 0 for automatic width
	}
	arrowsBarConfig.SetStyleUnicodeArrows()
	// Go
	if err := liveprogress.Start(); err != nil {
		panic(err)
	}
	hashRandom(size7G, liveprogress.DefaultConfig)
	hashRandom(size8G, arrowsBarConfig)
	hashRandom(size3G, liveprogress.DefaultConfig)
	liveprogress.AddCustomLine(spinner.Next)
	// Wait
	workers.Wait()
	liveprogress.Stop(true)
}

func hashRandom(size int, config liveprogress.BarConfig) {
	// Open random
	fd, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}
	// Create the hasher
	hasher := New(fd, size)
	// Create the hasher progress bar
	bar := liveprogress.AddBar(uint64(size), config,
		liveprogress.DecoratorAddition{
			Decorator: func(bar *liveprogress.Bar) string {
				return fmt.Sprintf("Hashing %d bytes ", size)
			},
			Prepend: true,
		},
		liveprogress.AppendPercent(),
		liveprogress.DecoratorAddition{
			Decorator: func(bar *liveprogress.Bar) string {
				return fmt.Sprintf("  SHA256: 0x%X", hasher.GetCurrentHash())
			},
		},
		liveprogress.DecoratorAddition{
			Decorator: func(bar *liveprogress.Bar) string {
				return "  Remaining:"
			},
		},
		liveprogress.AppendTimeRemaining(),
	)
	if bar == nil {
		panic("failed to create progress bar")
	}
	// Start hashing
	workers.Add(1)
	go func() {
		if err = hasher.ComputeHash(bar.CurrentAdd); err != nil {
			panic(err)
		}
		// Hashing done
		fmt.Fprintf(liveprogress.Bypass(), "%d bytes SHA256 done: 0x%X\n", size, hasher.GetCurrentHash())
		liveprogress.RemoveBar(bar)
		// Cleanup
		fd.Close()
		workers.Done()
	}()
}
