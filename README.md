# liveprogress
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hekmon/liveprogress)](https://pkg.go.dev/github.com/hekmon/liveprogress)

liveprogress is a golang library allowing to print and update progress bars on a terminal. It is heavily inspired by [uiprogress](https://github.com/gosuri/uiprogress) but redone on top of the forked [liveterm](https://github.com/hekmon/liveterm) library in order to take advantage of its enhancements.

## Examples

### Simple

Code available [here](examples/simple/main.go).

```go
// Config
liveprogress.DefaultConfig.Width = 70 // leave it a 0 for automatic width
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
```

TODO gif

### Advanced

See full source code [here](examples/advanced/main.go)

TODO gif

## Installation

```bash
go get -v github.com/hekmon/liveprogress
```
