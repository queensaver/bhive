package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
  "bytes"
  "strings"

	"encoding/json"
	"github.com/wogri/bbox/structs/temperature"
	"github.com/yryz/ds18b20"
	"net/http"
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
    r := strings.Replace(a, ":", "", -1)
		return r, nil
	}
	return "", nil
}

func post(t temperature.Temperature) error {
  j, err := json.Marshal(t)
  if err != nil {
    return err
  }
	req, err := http.NewRequest("POST", *serverAddr, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200" {
    log.Println("Got http return code ", resp.Status)
  }
  return nil
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
			BHiveID:      mac,
			SensorID:    sensor}
		err = post(t)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}
