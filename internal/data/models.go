package data

import (
	"log/slog"

	"gorm.io/gorm"
)

type Models struct {
	Location *LocationModel
	Device   *DeviceModel
	Module   *ModuleModel
	Data     *DataModel

	ModuleModels *ModuleModels
}

type ModuleModels struct {
	DB                *gorm.DB
	LightController   *LightControllerModel
	LightSensor       *LightSensorModel
	PresenceDetector  *PresenceDetectorModel
	LuminositySensor  *LuminositySensorModel
	TemperatureSensor *TemperatureSensorModel
	ConsumptionSensor *ConsumptionSensorModel
}

func NewModels(db *gorm.DB, broker *Broker, logger *slog.Logger) Models {
	return Models{
		Location: &LocationModel{DB: db},
		Device:   &DeviceModel{DB: db, Broker: broker},
		Module:   &ModuleModel{DB: db, Broker: broker},
		Data:     &DataModel{DB: db, Broker: broker, Logger: logger},

		ModuleModels: &ModuleModels{
			DB:                db,
			LightController:   &LightControllerModel{DB: db, Broker: broker},
			LightSensor:       &LightSensorModel{DB: db, Broker: broker},
			PresenceDetector:  &PresenceDetectorModel{DB: db, Broker: broker},
			LuminositySensor:  &LuminositySensorModel{DB: db, Broker: broker},
			TemperatureSensor: &TemperatureSensorModel{DB: db, Broker: broker},
			ConsumptionSensor: &ConsumptionSensorModel{DB: db, Broker: broker},
		},
	}
}
