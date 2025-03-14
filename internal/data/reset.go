package data

import (
	"fmt"

	"gorm.io/gorm"
)

const RESET = "reset"

type Reset struct {
	gorm.Model
	DeviceID  string `gorm:"index"`
	Name      string
	BoolValue bool
}

func (r Reset) GetValue() any {
	return r.BoolValue
}

func (r Reset) GetName() string {
	return r.Name
}

func NewResetModule() (IModule, error) {

	resetModule := &Module{
		Name:  RESET,
		Value: "1",
	}

	iModule, err := resetModule.ToIModule()
	if err != nil {
		return nil, fmt.Errorf("error creating reset iModule")
	}

	return iModule, nil
}
