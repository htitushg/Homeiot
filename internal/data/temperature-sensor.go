package data

import (
	"strconv"
	
	"gorm.io/gorm"
)

const TEMPERATURE_SENSOR = "temperatureSensor"

type TemperatureSensor struct {
	gorm.Model
	DeviceID     string `gorm:"index"`
	Name         string
	ValueDegrees float64
}

func (t TemperatureSensor) GetValue() any {
	return t.ValueDegrees
}

func (t TemperatureSensor) GetName() string {
	return t.Name
}

type TemperatureSensorModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *TemperatureSensorModel) Set(channel string, value any) error {
	floatValue, err := ToFloat(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatFloat(floatValue, 'f', 2, 64))
	return nil
}
