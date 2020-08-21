package main

import (
	"fmt"
  "flag"

	"github.com/yryz/ds18b20"
)

var serverAddr= flag.String("server_addr", "http://machine.intranet.wogri.com:8333/temperature", "HTTP server port")
var debug = flag.Bool("debug", false, "debug mode")

func main() {
	flag.Parse()
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
		}
	}
}
