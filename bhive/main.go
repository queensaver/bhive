package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"os"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	// "github.com/queensaver/bhive/bhive/sound"
	"github.com/queensaver/bhive/bhive/temperature"
	"github.com/queensaver/packages/logger"
	scaleStruct "github.com/queensaver/packages/scale"
	// soundStruct "github.com/queensaver/packages/sound"
	temperastureStruct "github.com/queensaver/packages/temperature"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

var serverAddr = flag.String("server_addr", "http://192.168.233.1:8333", "HTTP server port")
var ramDisk = flag.String("ramdisk", "/home/pi/bOS", "loccation of ramdisk to store temporary files")
var measurements = flag.Int("num_weight_measurements", 5, "Number of scale measurements")
// var soundFile = flag.String("sound_file", "/home/pi/bOS/audio.wav", "File where to record sound samples to")
var scaleConfigFile = flag.String("scale_config_file", "/home/pi/.queensaver_scale_config", "Scale config file")
var scanCmd = flag.String("scan_command", "/home/pi/capture.sh", "Command to execute for a varroa scan.")

func getMacAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error("Could not get interfaces", "error", err)
		return "", err
	}
	// find wlan0 interface
	for _, i := range interfaces {
		if i.Name == "wlan0" {
			a := i.HardwareAddr.String()
			if a != "" {
				r := strings.Replace(a, ":", "", -1)
				return r, nil
			}
		}
	}
	return "", fmt.Errorf("could not find wlan0 interface")
}

func post(req *http.Request) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return (err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Info("Post got unexpected return code", "http_status", resp.Status)
	}
	return nil
}

func sendFlush() error {
	req, err := http.NewRequest("POST", *serverAddr+"/flush", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	return post(req)
}

func postWeight(w scaleStruct.Scale) error {
	j, err := json.Marshal(w)
	if err != nil {
		return err
	}
	logger.Debug("posting weight ", "json", string(j))
	req, err := http.NewRequest("POST", *serverAddr+"/scale", bytes.NewBuffer(j))
	if err != nil {
		return (err)
	}
	req.Header.Set("Content-Type", "application/json")
	return post(req)
}

func postTemperature(t temperastureStruct.Temperature) error {
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", *serverAddr+"/temperature", bytes.NewBuffer(j))
	if err != nil {
		return (err)
	}
	req.Header.Set("Content-Type", "application/json")
	logger.Debug("posting temperature", "json", string(j))

	return post(req)
}

/*
func postSound(s *soundStruct.Sound) error {
	j, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", *serverAddr+"/sound", bytes.NewBuffer(j))
	if err != nil {
		return (err)
	}
	req.Header.Set("Content-Type", "application/json")
	logger.Debug("posting sound", "json", string(j))

	return post(req)
}
*/

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
		logger.Fatal("could not get mac address", "error", err)
	}
	logger.Debug("MAC address", "mac", mac)
	t, err := temperature.GetTemperature(mac)
	if err != nil {
		logger.Error("Error getting temperature", "error", err)
	}
	t.Timestamp = time.Now().Unix()
	logger.Debug("Temperature", "temperature", t)
	postTemperature(*t)

	err = write_python()
	if err != nil {
		logger.Fatal("Error writing python script", "error", err)
	}
	/*
	c, err := config.GetBHiveConfig(*serverAddr + "/config")
	if err != nil {
		logger.Fatal("Error getting config", "error", err)
	}
	logger.Debug("Config received", "config", c)
	*/
	// read the json config from the config file in /home/pi/.queensaver_scale_config
	conf, err := os.ReadFile(*scaleConfigFile)

	if err != nil {
		logger.Fatal("Error reading scale config file", "error", err)
	}

	type ScaleConfig struct {
		Offset float64 `json:"offset"`
		Scale float64 `json:"scale"`
	}

	sc := ScaleConfig{}
	err = json.Unmarshal(conf, &sc)
	if err != nil {
		logger.Fatal("Error unmarshalling config file", "error", err)
	}

	var weights []float64
	for i := 0; i < *measurements; i++ {
		weight, err := executePython(sc.Scale, sc.Offset)
		if err != nil {
			logger.Fatal("Error executing python script: ", "error", err)
		}
		weights = append(weights, weight)
	}
	sort.Float64s(weights)
	medianPosition := len(weights) / 2
	weight := weights[medianPosition] // We ignore that an even number of measurements would not calculate the exact median value.
	logger.Debug("Weight", "weight", weight)
	postWeight(scaleStruct.Scale{Weight: weight,
		BhiveId: mac,
		Epoch:   time.Now().Unix()})

		/*
	if c.RecordSound {
		logger.Info("Recording sound")
		audioRecording, err := sound.RecordSound(mac, int(c.SoundRecordingDuration), *soundFile)
		if err != nil {
			logger.Error("Error recording sound", "error", err)
		}
		err = postSound(audioRecording)
		if err != nil {
			logger.Info("Error posting sound to bbox", "error", err)
		}
	}
	*/


	cmd := exec.Command(*scanCmd)
	err = cmd.Run()
  
	if err != nil {
	  logger.Error("Scan command failed", "error", err)
	}

	err = sendFlush()
	if err != nil {
		logger.Info("Error sending flush", "error", err)
	}
}
