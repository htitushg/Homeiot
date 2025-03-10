package data

import (
	"fmt"
	"strconv"
	
	"gorm.io/gorm"
)

const LIGHT_CONTROLLER = "lightController"

type LightController struct {
	gorm.Model
	DeviceID string `gorm:"index"`
	Name     string
	On       bool
}

func (l LightController) GetValue() any {
	return l.On
}

func (l LightController) GetName() string {
	return l.Name
}

func (l LightController) getChannel() (string, error) {
	return fmt.Sprintf("home/%s/%d/%s/%s/%s"), nil
}

type LightControllerModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *LightControllerModel) Set(channel string, value any) error {
	boolValue, err := ToBool(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatBool(boolValue))
	return nil
}
