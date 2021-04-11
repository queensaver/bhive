package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	scaleStruct "github.com/btelemetry/packages/scale"
	temperastureStruct "github.com/btelemetry/packages/temperature"
	"github.com/btelemetry/bhive/bhive/temperature"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

var serverAddr = flag.String("server_addr", "http://192.168.233.1:8333", "HTTP server port")
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

func post(req *http.Request) error {
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

func postWeight(w scaleStruct.Scale) error {
	j, err := json.Marshal(w)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", *serverAddr+"/scale", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")
	return post(req)
}

func postTemperature(t temperastureStruct.Temperature) error {
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", *serverAddr+"/temperature", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	return post(req)
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
	err = ioutil.WriteFile(*ramDisk+"/hx711.py", pyHx711, 0755)
	if err != nil {
		return err
	}
	return nil
}

func execute_python() (float64, error) {
	var err error
	cmd := exec.Command("python3", *ramDisk+"/scale.py") // TODO: add calibration parameters
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	buf, err := io.ReadAll(stdout)
	if err != nil {
		return 0, err
	}

	stderrBuf, err := io.ReadAll(stderr)
	if err != nil {
		return 0, err
	}

	if err := cmd.Wait(); err != nil {
		return 0, err
	}
	fmt.Println("Python StdErr Output: ", string(stderrBuf))
	weight, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		return 0, err
	}
	return weight, nil
}

func main() {
	flag.Parse()
	mac, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
	}
	t, err := temperature.GetTemperature(mac)
	if err != nil {
		log.Println("Error getting temperature: ", err)
	} else {
		t.Timestamp = time.Now().Unix()
		fmt.Println("Temperature: ", t)
		postTemperature(*t)
	}

	err = write_python()
	if err != nil {
		log.Println("Error writing python script: ", err)
	}
	weight, err := execute_python()
	if err != nil {
		log.Println("Error executing python script: ", err)
	} else {
		fmt.Println("Weight: %s", weight)
		postWeight(scaleStruct.Scale{Weight: weight,
			BHiveID:   mac,
			Timestamp: time.Now().Unix()})
	}
}
