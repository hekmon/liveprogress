package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveprogress"
)

var (
	size1G = 1 * 1024 * 1024 * 1024
	size3G = 3 * 1024 * 1024 * 1024
	size7G = 5 * 1024 * 1024 * 1024
)

func main() {
	// Config
	liveprogress.SetProgressStyleUTF8Arrows()
	var workers sync.WaitGroup
	// Go
	liveprogress.Start()
	workers.Add(1)
	go func() {
		hashRandom(size3G)
		workers.Done()
	}()
	workers.Add(1)
	go func() {
		hashRandom(size7G)
		workers.Done()
	}()
	workers.Add(1)
	go func() {
		hashRandom(size1G)
		workers.Done()
	}()
	// Wait
	time.Sleep(10 * time.Millisecond)
	var spinner liveprogress.Spinner
	liveprogress.AddCustomLine(spinner.Next)
	workers.Wait()
	liveprogress.Stop(true)
}

func hashRandom(size int) {
	// Open random
	fd, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	// Create the hasher
	hasher := New(fd, size)
	// Create the hasher progress bar
	bar := liveprogress.AddBar(uint64(size))
	bar.PrependFunc(func(bar *liveprogress.Bar) string {
		return fmt.Sprintf("Hashing %d bytes ", size)
	})
	bar.AppendPercent()
	bar.AppendFunc(func(bar *liveprogress.Bar) string {
		return "  Remaining:"
	})
	bar.AppendTimeRemaining()
	bar.AppendFunc(func(bar *liveprogress.Bar) string {
		return fmt.Sprintf("  SHA256: %X", hasher.GetCurrentHash())
	})
	// Start hashing
	err = hasher.ComputeHash(bar.CurrentAdd)
	if err != nil {
		panic(err)
	}
	// Hashing done
	liveprogress.RemoveBar(bar)
	fmt.Fprintf(liveprogress.Bypass(), "%d bytes file hash: %X\n", size, hasher.GetCurrentHash())
}
