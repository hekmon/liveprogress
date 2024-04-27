package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveprogress/v2"
	"github.com/hekmon/liveprogress/v2/colors"
)

const (
	// sizes
	size3G = 3 * 1024 * 1024 * 1024
	size5G = 5 * 1024 * 1024 * 1024
	size8G = 8 * 1024 * 1024 * 1024
	// bar width for example
	barWidth = 40
)

var (
	workers sync.WaitGroup
	spinner liveprogress.Spinner
)

func main() {
	// Global config (these are already the default values)
	liveprogress.Output = os.Stdout
	liveprogress.RefreshInterval = 100 * time.Millisecond
	// Go
	if err := liveprogress.Start(); err != nil {
		panic(err)
	}
	// Main line, let's spin while hashing
	liveprogress.SetMainLineAsCustomLine(spinner.Next)
	// Create some colored styles
	colors.Generate() // Call it after Start() if you have changed the default Output value, otherwise you can omit it
	basicANSIGreenColor := colors.ANSIBasicGreen
	extendedAINSIPurpleColor := colors.ANSIExtended93
	rgbPinkColor := colors.RGB("#ff5faf")
	// File 1
	hashRandom(size5G,
		liveprogress.WithWidth(barWidth),
		liveprogress.WithPlainRunes(),
		liveprogress.WithBarStyle(rgbPinkColor),
		liveprogress.WithAppendPercent(rgbPinkColor.Bold()),
	)
	// File 2
	hashRandom(size8G,
		liveprogress.WithWidth(barWidth),
		liveprogress.WithLineFillRunes(),
		liveprogress.WithBarStyle(extendedAINSIPurpleColor),
		liveprogress.WithAppendPercent(extendedAINSIPurpleColor.Bold()),
	)
	// File 3
	hashRandom(size3G,
		liveprogress.WithWidth(barWidth),
		liveprogress.WithLineBracketsRunes(),
		liveprogress.WithBarStyle(basicANSIGreenColor),
		liveprogress.WithAppendPercent(basicANSIGreenColor.Bold()),
	)
	// Wait
	workers.Wait()
	if err := liveprogress.Stop(true); err != nil {
		panic(err)
	}
}

func hashRandom(size int, opts ...liveprogress.BarOption) {
	// Open random
	fd, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}
	// Create the hasher
	hasher := New(fd, size)
	// default options
	bold := colors.NoColor.Bold()
	defaultOpts := []liveprogress.BarOption{
		liveprogress.WithTotal(uint64(size)),
		liveprogress.WithPrependDecorator(func(bar *liveprogress.Bar) string {
			return fmt.Sprintf("Hashing %d bytes ", size)
		}),
		liveprogress.WithAppendDecorator(func(bar *liveprogress.Bar) string {
			return fmt.Sprintf("  SHA256: %s", bold.Styled(fmt.Sprintf("0x%X", hasher.GetCurrentHash())))
		}),
		liveprogress.WithAppendDecorator(func(bar *liveprogress.Bar) string {
			return "  Remaining:"
		}),
		liveprogress.WithAppendTimeRemaining(colors.NoColor),
	}

	// Create the hasher progress bar
	bar := liveprogress.AddBar(append(opts, defaultOpts...)...)
	if bar == nil {
		panic("failed to create progress bar")
	}
	// Start hashing
	workers.Add(1)
	go func() {
		defer workers.Done()
		defer fd.Close()
		if err = hasher.ComputeHash(bar.CurrentAdd); err != nil {
			panic(err)
		}
		liveprogress.RemoveBar(bar)
		fmt.Fprintf(liveprogress.Bypass(), "%d bytes SHA256 done: 0x%X\n", size, hasher.GetCurrentHash())
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
