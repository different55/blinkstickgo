/* 
 * BlinkStickGo - A libusb-based go package for controlling the BlinkStick line of products.
 * 
 *   This Source Code Form is subject to the terms of the Mozilla Public
 *   License, v. 2.0. If a copy of the MPL was not distributed with this
 *   file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * blinkstick_test.go
 */

package blinkstickgo

import "testing"

// Basic usage, setting all LEDs white
func ExampleBlinkStick() {
	Init()
	defer Fini()

	sticks, err := FindAll()
	if err != nil {
		panic(err)
	} else if len(sticks) == 0 {
		panic("No connected BlinkStick devices for testing")
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
}

func TestSetRGB(t *testing.T) {
	Init()
	defer Fini()
	
	sticks, err := FindAll()
	if err != nil {
		panic(err)
	} else if len(sticks) == 0 {
		panic("No connected BlinkStick devices for testing")
	}

	for _, stick := range sticks {
		err := stick.SetRGB(0, 0, 255, 255, 255)
		if err != nil {
			panic(err)
		}

		recvData, err := stick.GetLEDData(1)
		for _, chunk := range recvData {
			if chunk < 252 { // Can't check if it == 255 because the Blinkstick chews things up a bit
				t.Fail()
			}
		}
	}
}