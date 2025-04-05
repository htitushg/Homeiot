package data

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Broker struct {
	opts mqtt.ClientOptions
	mqtt.Client
	qos byte
}

func NewBroker(host string, port int64, qos byte) *Broker {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", host, port))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Broker{
		opts:   *opts,
		Client: client,
		qos:    qos,
	}
}

func (b *Broker) Pub(topic, message string) {
	token := b.Publish(topic, b.qos, false, message)
	token.Wait()
}
