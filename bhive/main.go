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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/queensaver/bhive/bhive/temperature"
	"github.com/queensaver/packages/config"
	scaleStruct "github.com/queensaver/packages/scale"
	temperastureStruct "github.com/queensaver/packages/temperature"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

var serverAddr = flag.String("server_addr", "http://192.168.233.1:8333", "HTTP server port")
var ramDisk = flag.String("ramdisk", "/home/pi/bOS", "loccation of ramdisk to store temporary files")
var measurements = flag.Int("num_weight_measurements", 5, "Number of scale measurements")

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
	log.Println("posting weight ", string(j))
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
	log.Println("posting temperature", string(j))

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

func executePython(reference_unit float64, offset float64) (float64, error) {
	var err error
	cmd := exec.Command("python3",
		*ramDisk+"/scale.py",
		fmt.Sprintf("--reference_unit=%f", reference_unit),
		fmt.Sprintf("--offset=%f", offset))
	fmt.Println("executing ", cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		stderrBuf, err := io.ReadAll(stderr)
		if err != nil {
			return 0, err
		}
		fmt.Println("Python StdErr Output: ", string(stderrBuf))
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
	log.Println("MAC: ", mac)
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
		log.Fatalln("Error writing python script: ", err)
	}
	c, err := config.GetBHiveConfig(*serverAddr + "/config")
	if err != nil {
		log.Fatalln("Error getting config: ", err)
	}
	log.Println("Config received: ", c)
	var weights []float64
	for i := 0; i < *measurements; i++ {
		weight, err := executePython(c.ScaleReferenceUnit, c.ScaleOffset)
		if err != nil {
			log.Fatalln("Error executing python script: ", err)
		}
		weights = append(weights, weight)
	}
	sort.Float64s(weights)
	medianPosition := len(weights) / 2
	weight := weights[medianPosition] // We ignore that an even number of measurements would not calculate the exact median value.
	postWeight(scaleStruct.Scale{Weight: weight,
		BhiveId: mac,
		Epoch:   time.Now().Unix()})
}
