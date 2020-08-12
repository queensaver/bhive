package main

import (
	"flag"
	"fmt"
	//"io/ioutil"

	"github.com/MichaelS11/go-hx711"
	"github.com/wogri/bhive/scale/thingspeak_client"
)

var thingspeakKey = flag.String("thingspeak_api_key", "48PCU5CAQ0BSP4CL", "API key for Thingspeak")

func main() {
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

	defer hx711.Shutdown()
	err = hx711.Reset()
	if err != nil {
		fmt.Println("Reset error:", err)
		return
	} // SetGain default is 128
	// Gain of 128 or 64 is input channel A, gain of 32 is input channel B
	//hx711.SetGain(128)

	// make sure to use your values from calibration above
	hx711.AdjustZero = 43428
	hx711.AdjustScale = 20.544371

	// previousReadings := []float64{}
	// movingAvg, err := hx711.ReadDataMedianThenMovingAvgs(11, 8, &previousReadings)
	weight, err := hx711.ReadDataMedian(11)
	if err != nil {
		fmt.Println("ReadDataMedianThenMovingAvgs error:", err)
	}
	fmt.Println(weight)
	thing := thingspeak_client.NewChannelWriter(*thingspeakKey)
	// avg := fmt.Sprintf("%f", movingAvg)
	// fmt.Println(avg)
	thing.AddField(1, weight)
	_, err = thing.Update()
	if err != nil {
		fmt.Println("ThingSpeak error:", err)
	}
	//fmt.Println("HTTP: %s", r.Status)
	//body, err := ioutil.ReadAll(r.Body)
	//fmt.Println(string(body))

}
