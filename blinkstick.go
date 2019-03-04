// Package blinkstickgo provides functions to interact with the BlinkStick line of products.
package blinkstickgo

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/google/gousb"
)

var ctx *gousb.Context

// Init initializes the USB library
func Init() {
	ctx = gousb.NewContext()
}

// Fini closes the USB context used
func Fini() {
	ctx.Close()
}

const vendorID = 0x20A0
const productID = 0x41E5

// FindAll detects and returns all BlinkSticks connected to the system
func FindAll() ([]BlinkStick, error) {
	var blinksticks []BlinkStick

	devices, err := ctx.OpenDevices(filterBlinkStick)
	if err != nil {
		return blinksticks, err
	}

	for _, device := range devices {
		serial, err := device.SerialNumber()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not grab Serial for BlinkStick device", err)
		}
		blinksticks = append(blinksticks, BlinkStick{device, false, serial, 0})
	}
	return blinksticks, nil
}

// The BlinkStick struct represents an individual BlinkStick device.
type BlinkStick struct {
	device   *gousb.Device
	Inverse  bool
	Serial   string
	ledCount int
}

// GetLEDCount returns the number of LEDs for supported devices
func (stk *BlinkStick) GetLEDCount() int {
	if stk.ledCount == 0 {
		var buffer = make([]byte, 2)

		responseLen, err := stk.device.Control(0x80|0x20, 0x01, 0x81, 0x00, buffer)
		if err != nil || responseLen < 2 {
			return -1
		}

		stk.ledCount = int(buffer[1])
	}

	return stk.ledCount
}

// GetName returns the name of the device
func (stk *BlinkStick) GetName() string {
	var buffer = make([]byte, 33)

	err := stk.control(0x80|0x20, 0x01, 0x02, 0x00, buffer)
	if err != nil {
		return ""
	}

	buffer = append(buffer, 0) // Just in case there's no terminating null byte, we'll add our own.
	return string(buffer)
}

// GetInfo returns the name of the device
func (stk *BlinkStick) GetInfo() string {
	var buffer = make([]byte, 33)

	err := stk.control(0x80|0x20, 0x01, 0x03, 0x00, buffer)
	if err != nil {
		return ""
	}

	buffer = append(buffer, 0) // Just in case there's no terminating null byte, we'll add our own.
	return string(buffer)
}

// SetName writes a new name for the device to info block one.
// If you're worried about extreme longevity, use sparingly. I hear this stuff
// can only withstand so many writes.
func (stk *BlinkStick) SetName(name string) error {
	return stk.control(0x20, 0x09, 0x02, 0x00, []byte(name))
}

// SetInfo writes a new block of data to info block two.
// If you're worried about extreme longevity, use sparingly. I hear this stuff
// can only withstand so many writes.
func (stk *BlinkStick) SetInfo(info string) error {
	return stk.control(0x20, 0x09, 0x03, 0x00, []byte(info))
}

// SetRGB sends a color to the device in RGB format
func (stk *BlinkStick) SetRGB(channel, index, r, g, b byte) error {
	if stk.Inverse {
		r, g, b = 255-r, 255-g, 255-b
	}

	if index == 0 && channel == 0 {
		return stk.control(0x20, 0x09, 0x01, 0x00, []byte{0, r, g, b})
	}
	return stk.control(0x20, 0x09, 0x01, 0x00, []byte{0, r, g, b})
}

// SetRandom sends a random color to the device
func (stk *BlinkStick) SetRandom(channel, index byte) error {
	var rand = rand.Uint32()
	return stk.SetRGB(channel, index, byte(rand>>24), byte(rand>>16), byte(rand>>8))
}

// SetLEDData updates the entire stick with a slice of alternating RGB values
func (stk *BlinkStick) SetLEDData(channel byte, data []byte) error {
	var reportID, maxLEDs = stk.getReportID(len(data))
	var report = []byte{0, channel}

	for i := 0; uint16(i) < maxLEDs*3; i++ {
		if len(data) > i {
			report = append(report, data[i])
		} else {
			report = append(report, 0)
		}
	}
	return stk.control(0x20, 0x09, reportID, 0x00, report)
}

// A razor thin wrapper around gousb.Device.Control
func (stk *BlinkStick) control(requestType, request uint8, val, idx uint16, data []byte) error {
	_, err := stk.device.Control(requestType, request, val, idx, data)
	return err
}

// Returns true if the device is a BlinkStick
func filterBlinkStick(desc *gousb.DeviceDesc) bool {
	return desc.Vendor == vendorID && desc.Product == productID
}

// The BlinkStick seems to use different Report IDs for different data lengths when setting all LEDs.
func (stk *BlinkStick) getReportID(count int) (uint16, uint16) {
	var reportID uint16 = 9
	var maxLEDs uint16 = 64

	switch {
	case count <= 8*3:
		maxLEDs = 8
		reportID = 6
	case count <= 16*3:
		maxLEDs = 16
		reportID = 7
	case count <= 32*3:
		maxLEDs = 32
		reportID = 8
	case count <= 64*3:
		maxLEDs = 64
		reportID = 9
	}

	return reportID, maxLEDs
}
