package data

import (
	"strconv"
	
	"gorm.io/gorm"
)

const LIGHT_SENSOR = "lightSensor"

type LightSensor struct {
	gorm.Model
	DeviceID string `gorm:"index"`
	Name     string
	IsOn     bool
}

func (l LightSensor) GetValue() any {
	return l.IsOn
}

func (l LightSensor) GetName() string {
	return l.Name
}

type LightSensorModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *LightSensorModel) Set(channel string, value any) error {
	boolValue, err := ToBool(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatBool(boolValue))
	return nil
}
