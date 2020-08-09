package main

import (
	"fmt"

	"github.com/MichaelS11/go-hx711"
)

func main() {
	err := hx711.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	// hx711, err := hx711.NewHx711("GPIO6", "GPIO5")
	hx711, err := hx711.NewHx711("6", "5")
	if err != nil {
		fmt.Println("NewHx711 error:", err)
		return
	}

	// SetGain default is 128
	// Gain of 128 or 64 is input channel A, gain of 32 is input channel B
	hx711.SetGain(128)

	// make sure to use your values from calibration above
	hx711.AdjustZero = 43428
	hx711.AdjustScale = 20.544371

	previousReadings := []float64{}
	movingAvg, err := hx711.ReadDataMedianThenMovingAvgs(11, 8, &previousReadings)
	if err != nil {
		fmt.Println("ReadDataMedianThenMovingAvgs error:", err)
	}

	// moving average
	fmt.Println(movingAvg)
	// fmt.Println(previousReadings)

}
