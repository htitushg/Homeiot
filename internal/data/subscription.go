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

/**
 * Sub subscribes to the given topic and sets the appropriate handler.
 * The handler is determined based on the topic prefix and suffix.
 * - If the topic starts with "home/" and ends with "/startup", the startupHandler is used.
 */
func (m *DataModel) Sub(topic string) {
	// Ajouté le 4/04/2025 à 10h06
	// DEBUG
	m.Logger.Debug("Sub subscribing to MQTT topic", slog.String("TOPIC", topic))

	token := m.Broker.Subscribe(topic, m.Broker.qos, m.mqttHandler)
	token.Wait()
}

func (m *DataModel) mqttHandler(client mqtt.Client, msg mqtt.Message) {
	if strings.HasPrefix(msg.Topic(), "home/") {
		if strings.HasSuffix(msg.Topic(), "/startup") {
			// DEBUG
			m.Logger.Debug("SUB Dans la boucle subscribing to MQTT startup topic", slog.String("TOPIC", msg.Topic()))
			m.startupHandler(client, msg)
		} else {
			m.dataHandler(client, msg)
		}
	} else {
		m.messageHandler(client, msg)
	}
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

/**
 * startupHandler handles the startup message from the device.
 * It checks if the device exists in the database and creates it if not.
 */
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
	// DEBUG
	m.Logger.Debug("startupHandler Check if the Device exists and create it if not", slog.String("HANDLER", "messageHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):

			// Create the Device
			// DEBUG
			m.Logger.Debug("startupHandler Create the Device", slog.String("HANDLER", "messageHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))

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
	// DEBUG
	m.Logger.Debug("startupHandler Répondre à l'appareil avec les données récupérées ou créées", slog.String("JSON MESSAGE", string(jsonMessage)))
	// DEBUG
	m.Logger.Debug("startupHandler Valeur du Topic", slog.String("TOPIC", msg.Topic()))

	// Il faut peut être modifier le startup en setup dans le msg.Topic
	// msg.Topic()= strings.Replace(msg.Topic(), "startup", "setup", 1)
	m.Broker.Pub(device.GetChannel(&Setup{}), string(jsonMessage))
}

func (m *DataModel) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// FIXME -> remove or modify to accommodate normal usage!
	// LOG WARNING MESSAGE
	m.Logger.Warn("received unknown MQTT message", slog.String("HANDLER", "messageHandler"), slog.String("TOPIC", msg.Topic()), slog.String("PAYLOAD", string(msg.Payload())))
}
