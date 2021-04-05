package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	temperastureStruct "github.com/wogri/bbox/packages/temperature"
	"github.com/wogri/bhive/packages/temperature"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

var serverAddr = flag.String("server_addr", "http://192.168.233.1:8333/temperature", "HTTP server port")

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

func post(t temperastureStruct.Temperature) error {
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", *serverAddr, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return (err)
	}
	defer resp.Body.Close()

	if resp.Status != "200" {
		log.Println("Got http return code ", resp.Status)
	}
	return nil
}

func main() {
	mac, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
	}
	t, err := temperature.GetTemperature(mac)
	if err != nil {
		log.Println("Error getting temperature: ", err)
	}
	post(*t)
	print(string(pyScale))
	print(string(pyHx711))
}
