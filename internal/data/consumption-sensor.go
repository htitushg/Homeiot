package data

import (
	"strconv"
	
	"gorm.io/gorm"
)

const CONSUMPTION_SENSOR = "consumptionSensor"

type ConsumptionSensor struct {
	gorm.Model
	DeviceID      string `gorm:"index"`
	Name          string
	ValueWattHour float64
}

func (c ConsumptionSensor) GetValue() any {
	return c.ValueWattHour
}

func (c ConsumptionSensor) GetName() string {
	return c.Name
}

type ConsumptionSensorModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *ConsumptionSensorModel) Set(channel string, value any) error {
	floatValue, err := ToFloat(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatFloat(floatValue, 'f', 2, 64))
	return nil
}
