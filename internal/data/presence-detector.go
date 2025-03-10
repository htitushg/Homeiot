package data

import (
	"strconv"
	
	"gorm.io/gorm"
)

const PRESENCE_DETECTOR = "presenceDetector"

type PresenceDetector struct {
	gorm.Model
	DeviceID   string `gorm:"index"`
	Name       string
	IsPresence bool
}

func (p PresenceDetector) GetValue() any {
	return p.IsPresence
}

func (p PresenceDetector) GetName() string {
	return p.Name
}

type PresenceDetectorModel struct {
	DB     *gorm.DB
	Broker *Broker
}

func (m *PresenceDetectorModel) Set(channel string, value any) error {
	boolValue, err := ToBool(value)
	if err != nil {
		return err
	}
	m.Broker.Pub(channel, strconv.FormatBool(boolValue))
	return nil
}
