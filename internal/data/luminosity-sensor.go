package data

import (
	"strconv"
	
	"gorm.io/gorm"
)

const LUMINOSITY_SENSOR = "luminositySensor"

type LuminositySensor struct {
	gorm.Model
	DeviceID   string `gorm:"index"`
	Name       string
	ValueLumen float64
}

func (l LuminositySensor) GetValue() any {
	return l.ValueLumen
}

func (l LuminositySensor) GetName() string {
	return l.Name
}

type LuminositySensorModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *LuminositySensorModel) Set(channel string, value any) error {
	floatValue, err := ToFloat(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatFloat(floatValue, 'f', 2, 64))
	return nil
}
