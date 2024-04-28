# liveprogress
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hekmon/liveprogress)](https://pkg.go.dev/github.com/hekmon/liveprogress)

liveprogress is a golang library allowing to print and update progress bars on a terminal. It is heavily inspired by [uiprogress](https://github.com/gosuri/uiprogress) but redone on top of the forked [liveterm](https://github.com/hekmon/liveterm/v2) library in order to take advantage of its enhancements.

In addition of the features of [liveterm](https://github.com/hekmon/liveterm/v2), it also add (or changes):
* Automatic bar length if its `width` is 0
* Bars characters are runes (Unicode support thru `liveterm`)
* Remove unecessary mutexes
	* usage of atomic operations for bar progress
	* decorators can be added only when instanciating the bar
* Custom (dynamic) lines that can be anything (not necessarly a progress bar)
* Main line concept: a bar or a custom line that will always be printed last (usefull for global progress when others lines indicate specific progress)
* Ability to style the bar and decorators using [termenv](https://github.com/muesli/termenv) styles

## Examples

### Simple

Code available [here](examples/simple/main.go).

```go
if err := liveprogress.Start(); err != nil {
	panic(err)
}
bar := liveprogress.AddBar(
	liveprogress.WithWidth(75), // remove for automatic size
	liveprogress.WithPrependPercent(liveprogress.BaseStyle()),
)
// By default a bar total is set to 100
for i := 0; i < liveprogress.DefaultTotal; i++ {
	// Wait a random time
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	// Increment the bar
	bar.CurrentIncrement()
}
if err := liveprogress.Stop(true); err != nil {
	panic(err)
}
fmt.Println("By setting the Stop() bool parameter to true, the progress bar is cleared at stop.")
```

![Simple example output animation](https://github.com/hekmon/liveprogress/blob/main/examples/simple/example.gif?raw=true)

### Advanced

See full source code [here](examples/advanced/main.go).

![Advanced example output animation](https://github.com/hekmon/liveprogress/blob/main/examples/advanced/example.gif?raw=true)

## Installation

```bash
go get -v github.com/hekmon/liveprogress/v2
```
