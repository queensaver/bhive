package temperature

import (
	"fmt"

	t "github.com/queensaver/packages/temperature"
	"github.com/queensaver/openapi/golang/proto/models"
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
		return &t.Temperature{
      Temperature: models.Temperature{BhiveId: mac},
			Error:   fmt.Sprintf("%s", err)}, err
	}

	for _, sensor := range sensors {
		measured_temperature, err := ds18b20.Temperature(sensor)
		if err != nil {
			return &t.Temperature{
        Temperature: models.Temperature{BhiveId: mac},
				Error:   fmt.Sprintf("%s", err)}, err
		}
		return &t.Temperature{
      Temperature: models.Temperature{Temperature: float32(measured_temperature),
			  BhiveId: mac}}, nil
	}
	return &t.Temperature{
		Temperature: models.Temperature{BhiveId: mac},
		Error:   "no sensor found"}, nil
}
