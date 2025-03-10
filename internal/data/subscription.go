package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

func (m *DataModel) Sub(topic string) {
	var callback mqtt.MessageHandler
	if strings.HasPrefix(topic, "home/") {
		if strings.HasSuffix(topic, "/startup") {
			callback = m.startupHandler
		} else {
			callback = m.dataHandler
		}
	} else {
		callback = m.messageHandler
	}
	token := m.Broker.Subscribe(topic, m.Broker.qos, callback)
	token.Wait()
}

func (m *DataModel) dataHandler(client mqtt.Client, msg mqtt.Message) {
	// DEBUG
	m.Logger.Debug("received MQTT message", slog.String("HANDLER", "dataHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))
	
	data, err := m.NewData(msg)
	if err != nil {
		m.Logger.Error(fmt.Errorf("error creating data from MQTT message: %w", err).Error())
		m.Logger.Warn("aborting data creation")
		return
	}
	err = m.insert(data)
	if err != nil {
		m.Logger.Error(fmt.Errorf("error inserting data: %w", err).Error())
		m.Logger.Warn("aborting data creation")
	}
}

func (m *DataModel) startupHandler(client mqtt.Client, msg mqtt.Message) {
	// DEBUG
	m.Logger.Debug("received startup MQTT message", slog.String("HANDLER", "messageHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))
	
	// Parse the payload into a StartupMessage
	startupMessage, err := NewStartupMessage(msg.Payload())
	if err != nil {
		m.Logger.Error(err.Error())
		return
	}
	
	// Convert the StartupMessage into a Device
	device := startupMessage.ToDevice()
	
	// Check if the Device exists and create it if not
	err = m.Check(device)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			
			// Create the Device
			result := m.DB.Create(&device)
			if result.Error != nil {
				m.Logger.Error(fmt.Errorf("error creating the device: %w", result.Error).Error())
				return
			}
			if result.RowsAffected == 0 {
				m.Logger.Error(fmt.Errorf("error creating the device: %w", result.Error).Error())
				return
			}
		default:
			m.Logger.Error(err.Error())
			return
		}
		// TODO -> what to do after creating the device, if necessary
	}
	
	// Create new StartupMessage from device fetched or created
	responseMessage := NewResponseMessage(device)
	jsonMessage, err := json.Marshal(responseMessage)
	if err != nil {
		m.Logger.Error(fmt.Errorf("error marshaling json: %w", err).Error())
	}
	
	// Respond to the device with the data fetched or created
	m.Broker.Pub(msg.Topic(), string(jsonMessage))
}

func (m *DataModel) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// FIXME -> remove or modify to accommodate normal usage!
	// LOG WARNING MESSAGE
	m.Logger.Warn("received unknown MQTT message", slog.String("HANDLER", "messageHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))
}
