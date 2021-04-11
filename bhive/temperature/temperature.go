package temperature

import (
	"fmt"

	t "github.com/btelemetry/packages/temperature"
	"github.com/yryz/ds18b20"
)

func GetTemperature(mac string) (*t.Temperature, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GetTemperature:", err)
		}
	}()
	sensors, err := ds18b20.Sensors()
	if err != nil {
		return nil, err
	}

	for _, sensor := range sensors {
		measured_temperature, err := ds18b20.Temperature(sensor)
		if err != nil {
			return nil, err
		}
		return &t.Temperature{
			Temperature: measured_temperature,
			BHiveID:     mac,
			SensorID:    sensor}, nil
	}
	return nil, nil
}
