package main

import "github.com/rajmaniar/hx711"
import (
	"fmt"
	"time"
)

func main() {

	clock := "gpio6"
	data := "gpio5"

	h, err := hx711.New(data, clock)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	for err == nil {
		var data int32
		data, err = h.ReadData()
		fmt.Printf("Read from HX711: %v\n", data)
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Printf("Stopped reading because of: %v\n", err)
}
