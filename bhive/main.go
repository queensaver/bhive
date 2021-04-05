package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"

	temperastureStruct "github.com/wogri/bbox/packages/temperature"
	"github.com/wogri/bhive/packages/temperature"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

var serverAddr = flag.String("server_addr", "http://192.168.233.1:8333/temperature", "HTTP server port")
var ramDisk = flag.String("ramdisk", "/home/pi/bOS", "loccation of ramdisk to store temporary files")

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

// writes out a python file to the ramdisk.
// this is so that we can run the python script afterwards.
// The reason we do this in python is because for whatever reason python produces more accurate results when measuring the scale.
func write_python() error {
	var err error
	err = ioutil.WriteFile(*ramDisk+"/scale.py", pyScale, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(*ramDisk+"/hx711.py", pyhx711, 0755)
	if err != nil {
		return err
	}
	return nil
}

func execute_python() error {
	var err error
	err = exec.Command("python3", *ramDisk+"/scale.py") // TODO: add calibration parameters
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

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
	// TODO: implement retry logic
	post(*t)
	err = write_python()
	if err != nil {
		log.Println("Error writing python script: ", err)
	}
	err = execute_python()
	if err != nil {
		log.Println("Error executing python script: ", err)
	}
	print(string(pyScale))
	print(string(pyHx711))
}
