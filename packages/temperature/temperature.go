package temperature

import (
	t "github.com/wogri/bbox/packages/temperature"
	"github.com/yryz/ds18b20"
)

func GetTemperature(mac string) (*t.Temperature, error) {
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
			BBoxID:      mac,
			SensorID:    sensor}, nil
	}
	return nil, nil
}
