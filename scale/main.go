package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	//"io/ioutil"

	"github.com/MichaelS11/go-hx711"
	"github.com/wogri/bhive/scale/thingspeak_client"
)

var thingspeakKey = flag.String("thingspeak_api_key", "48PCU5CAQ0BSP4CL", "API key for Thingspeak")
var debug = flag.Bool("debug", false, "debug mode")
var thingspeakActive = flag.Bool("thingspeak", false, "Activate thingspeak API if set to true")

func main() {
	flag.Parse()

	err := hx711.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	hx711, err := hx711.NewHx711("GPIO6", "GPIO5")
	// hx711, err := hx711.NewHx711("6", "5")
	if err != nil {
		fmt.Println("NewHx711 error:", err)
		return
	}

	// SetGain default is 128
	// Gain of 128 or 64 is input channel A, gain of 32 is input channel B
	//hx711.SetGain(128)

	// make sure to use your values from calibration above
	hx711.AdjustZero = 43428
	hx711.AdjustScale = 20.544371

	// previousReadings := []float64{}
	// movingAvg, err := hx711.ReadDataMedianThenMovingAvgs(11, 8, &previousReadings)
	var weights []float64
	for i := 0; i < 11; i++ {

		err = hx711.Reset()
		if err != nil {
			log.Println("Reset error:", err)
			return
		}
		// hx711.waitForDataReady()
		rawWeight, err := hx711.ReadDataRaw()
		if err != nil {
			log.Println("RadDataRaw error:", err)
		}
		if *debug {
			log.Println("Raw Weight measured is: ", rawWeight)
		}
		err = hx711.Shutdown()
		if err != nil {
			log.Println("Shutdown error:", err)
		}
		if rawWeight == 65535 {
			i--
			if *debug {
				log.Println("Read error. Waiting 100ms.")
			}
			time.Sleep(100 * time.Millisecond)
			time.Sleep(1 * time.Second)

		}
		weight := float64(rawWeight-hx711.AdjustZero) / hx711.AdjustScale
		weights = append(weights, weight)
		if *debug {
			log.Println("measured weight is: ", weight)
		}
	}
	sort.Float64s(weights)
	medianWeight := weights[5]
	log.Println("median weight is: ", medianWeight)
	thing := thingspeak_client.NewChannelWriter(*thingspeakKey)
	// avg := fmt.Sprintf("%f", movingAvg)
	// fmt.Println(avg)
	thing.AddField(1, medianWeight)
	if *thingspeakActive {
		if *debug {
			log.Println("uploading data to Thingspeak...")
		}
		_, err = thing.Update()
		if err != nil {
			log.Println("ThingSpeak error:", err)
		}
	}
}
