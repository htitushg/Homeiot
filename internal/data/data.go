package data

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	DeviceID    string
	Device      Device `gorm:"foreignKey:DeviceID"`
	ModuleID    uint
	ModuleName  string
	ModuleValue string
}

type DataModel struct {
	DB     *gorm.DB
	Broker *Broker
	Logger *slog.Logger
}

func (m *DataModel) NewData(message mqtt.Message) (*Data, error) {

	// Parse channel name into single elements
	channel := message.Topic()
	channelElems := strings.Split(channel, "/")

	// Check if channel name follows normalized format
	if len(channelElems) != 6 {
		return nil, fmt.Errorf("invalid channel format")
	}

	// Get location type & ID
	locationType := channelElems[1]
	locationID, err := strconv.Atoi(channelElems[2])
	if err != nil {
		return nil, fmt.Errorf("error converting location ID %w", err)
	}

	// Get device type & ID
	device := channelElems[3]
	deviceID := channelElems[4]
	moduleName := channelElems[5]

	// Get value in payload
	moduleValue := string(message.Payload())
	if moduleValue == "" {
		return nil, fmt.Errorf("no value found in payload")
	}

	// Create data instance
	data := &Data{
		DeviceID: deviceID,
		Device: Device{
			ID:         deviceID,
			LocationID: uint(locationID),
			Location: Location{
				Model: gorm.Model{
					ID: uint(locationID),
				},
				Name: fmt.Sprintf("%s #%d", locationType, locationID),
				Type: locationType,
			},
			Name: fmt.Sprintf("%s #%s", device, deviceID),
		},
		ModuleName:  moduleName,
		ModuleValue: moduleValue,
	}

	// Retrieve device and module data from DB
	err = m.DB.Joins("Module").First(&data.Device, "id = ?", data.DeviceID).Error
	if err != nil {
		// FIXME -> reset device or skip data?
		return nil, fmt.Errorf("error finding device %w", err)
	}
	// Get ModuleID from Device.Modules by matching Device.ID/Module.Name
	for _, module := range data.Device.Modules {
		if module.Name == data.ModuleName {
			data.ModuleID = module.ID
		}
	}

	return data, nil
}

func (m *DataModel) insert(data *Data) error {
	result := m.DB.Model(&Data{}).Create(&data)
	if result.Error != nil {
		return fmt.Errorf("error inserting data: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("error inserting data: %d rows affected", result.RowsAffected)
	}
	return nil
}

func (m *DataModel) Check(device *Device) error {
	err := m.DB.Joins("Location").Joins("Module").First(&device, "id = ?", device.ID).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return err
		default:
			return fmt.Errorf("error fetching device %v: %w", device.ID, err)
		}
	}
	return nil
}
