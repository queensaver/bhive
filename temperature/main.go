package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"encoding/json"
	"github.com/wogri/bbox/structs/temperature"
	"github.com/yryz/ds18b20"
)

var serverAddr = flag.String("server_addr", "http://machine.intranet.wogri.com:8333/temperature", "HTTP server port")
var debug = flag.Bool("debug", false, "debug mode")

func getMacAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	a := interfaces[1].HardwareAddr.String()
	if a != "" {
		return a, nil
	}
	return "", nil
}

func main() {
	flag.Parse()
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	mac, err := getMacAddr()
	log.Println("mac: ", mac)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)

	for _, sensor := range sensors {
		measured_temperature, err := ds18b20.Temperature(sensor)
		if err != nil {
			log.Println(err)
			os.Exit(3)
		}
		t := temperature.Temperature{
			Temperature: measured_temperature,
			BBoxID:      mac,
			SensorID:    sensor}
		payload, err := json.Marshal(t)
		fmt.Printf(string(payload))
	}
}
