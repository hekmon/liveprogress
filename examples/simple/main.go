package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hekmon/liveprogress"
)

func main() {
	// Config
	liveprogress.DefaultConfig.Width = 78 // leave it a 0 for automatic width
	countTo := 100
	// Go
	if err := liveprogress.Start(); err != nil {
		panic(err)
	}
	bar := liveprogress.AddBar(uint64(countTo), liveprogress.DefaultConfig,
		liveprogress.PrependPercent(),
	)
	for i := 0; i < countTo; i++ {
		// Wait a random time
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		// Increment the bar
		bar.CurrentIncrement()
	}
	liveprogress.Stop(true)
	fmt.Println("By setting the Stop() bool parameter to true, the progress bar is cleared at stop.")
}
