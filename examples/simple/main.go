package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hekmon/liveprogress"
)

func main() {
	if err := liveprogress.Start(); err != nil {
		panic(err)
	}
	bar := liveprogress.AddBar(
		liveprogress.WithWidth(76),
		liveprogress.WithPrependPercent(""),
	)
	// By default a bar total is set to 100
	// As we omit it in AddBar()
	for i := 0; i < liveprogress.DefaultTotal; i++ {
		// Wait a random time
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		// Increment the bar
		bar.CurrentIncrement()
	}
	liveprogress.Stop(true)
	fmt.Println("By setting the Stop() bool parameter to true, the progress bar is cleared at stop.")
}
