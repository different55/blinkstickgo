# blinkstickgo
A libusb-based Go package for controlling the BlinkStick line of products

Currently, this package has only been tested on Linux, but theoretically it *might* work on your system, too.

# Installation
blinkstickgo only depends on gousb. But installing gousb requires that you have the headers for libusb on your system.

Linux users, it's pretty straightforward, just install your distro's libusb-dev package.  
For macOS, seems like you have [a few gotchas you may have to deal with](https://github.com/google/gousb#dependencies).  
Windows users, [you're going to need MINGW and I wish you luck](https://github.com/google/gousb#notes-for-installation-on-windows).

After you get that sorted out, just `go get github.com/different55/blinkstickgo`.

# Documentation
Documentation is available over on [godoc.org](https://godoc.org/github.com/different55/blinkstickgo).

# Basic Example
Turning all connected LEDs to white.
```go
blinkstickgo.Init()
defer blinkstickgo.Fini()

sticks, err := blinkstickgo.FindAll()
if err != nil {
	panic(err)
} else if len(sticks) == 0 {
	panic("No connected BlinkStick devices")
}

for _, stick := range sticks {
	count := stick.GetLEDCount()
	if count < 1 { // The Pros report their count as -1
		err := stick.SetRGB(0, 0, 255, 255, 255)
		if err != nil {
			panic(err)
		}
	} else {
		err := stick.SetAllRGB(0, 255, 255, 255)
		if err != nil {
			panic(err)
		}
	}
}
```
