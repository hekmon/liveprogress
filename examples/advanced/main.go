package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
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

type SHA256Progress struct {
	source    *readerCounter
	hasher    *hasherReporter
	hash      []byte
	computing sync.Mutex
}

func New(reader io.Reader, totalBytes int) *SHA256Progress {
	if reader == nil || totalBytes <= 0 {
		return nil
	}
	return &SHA256Progress{
		source: &readerCounter{
			source:  reader,
			maxRead: totalBytes,
		},
	}
}

func (hp *SHA256Progress) ComputeHash(reportWritten func(uint64)) (err error) {
	defer hp.computing.Unlock()
	hp.computing.Lock()
	if hp.hash != nil {
		return
	}
	// Prepare the hasher
	hp.hasher = &hasherReporter{
		dest:        sha256.New(),
		writeReport: reportWritten,
	}
	// Start copy
	if _, err = io.Copy(hp.hasher, hp.source); err != nil {
		return
	}
	hp.hash = hp.hasher.GetCurrentHash()
	hp.hasher = nil
	return
}

func (hp *SHA256Progress) GetCurrentHash() []byte {
	if hp.hash != nil {
		// compute is already done
		return hp.hash
	}
	if hp.hasher == nil {
		// too soon, CompupteHash has not been called yet
		return nil
	}
	return hp.hasher.GetCurrentHash()
}

type readerCounter struct {
	source  io.Reader
	read    int
	maxRead int
}

func (rc *readerCounter) Read(p []byte) (n int, err error) {
	if rc.read >= rc.maxRead {
		err = io.EOF
		return
	}
	n, err = rc.source.Read(p)
	rc.read += n
	return
}

type hasherReporter struct {
	dest        hash.Hash
	writeReport func(uint64)
	access      sync.Mutex
}

func (hc *hasherReporter) Write(p []byte) (n int, err error) {
	hc.access.Lock()
	n, err = hc.dest.Write(p)
	hc.writeReport(uint64(n))
	hc.access.Unlock()
	return
}

func (hc *hasherReporter) GetCurrentHash() (hash []byte) {
	hc.access.Lock()
	hash = hc.dest.Sum(nil)
	hc.access.Unlock()
	return
}
