package main

import (
	"fmt"
	"os"
	"sync"

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
	// Config
	liveprogress.SetProgressStyleUTF8Arrows()
	liveprogress.Start()
	// Go
	hashRandom(size7G)
	hashRandom(size8G)
	hashRandom(size3G)
	liveprogress.AddCustomLine(spinner.Next)
	// Wait
	workers.Wait()
	liveprogress.Stop(true)
}

func hashRandom(size int) {
	// Open random
	fd, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}
	// Create the hasher
	hasher := New(fd, size)
	// Create the hasher progress bar
	bar := liveprogress.AddBar(uint64(size))
	bar.PrependFunc(func(bar *liveprogress.Bar) string {
		return fmt.Sprintf("Hashing %d bytes ", size)
	})
	bar.AppendPercent()
	bar.AppendFunc(func(bar *liveprogress.Bar) string {
		return fmt.Sprintf("  SHA256: 0x%X", hasher.GetCurrentHash())
	})
	// Start hashing
	workers.Add(1)
	go func() {
		if err = hasher.ComputeHash(bar.CurrentAdd); err != nil {
			panic(err)
		}
		// Hashing done
		fmt.Fprintf(liveprogress.Bypass(), "%d bytes hashed: 0x%X\n", size, hasher.GetCurrentHash())
		liveprogress.RemoveBar(bar)
		// Cleanup
		fd.Close()
		workers.Done()
	}()
}
