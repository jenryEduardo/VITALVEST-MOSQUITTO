package main

import (
	"fmt"
	"log"
	"mosquitto/controllers"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://52.203.81.35:1883").
		SetClientID("go-mqtt-client")

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println(" Conectado al broker Mosquitto")

		// Suscribirse al tópico 'sensores/datos' con su handler
		if token := c.Subscribe("sensores/datos", 0, controllers.MessageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf(" Error al suscribirse a 'sensores/datos': %v", token.Error())
		}
		fmt.Println(" Suscrito al topic 'sensores/datos'")

		// // Suscribirse al tópico 'GSR-SENSOR' con su handler específico
		if token := c.Subscribe("GSR-SENSOR", 0, controllers.MessageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf(" Error al suscribirse a 'GSR-SENSOR': %v", token.Error())
		}
		fmt.Println(" Suscrito al topic 'GSR-SENSOR'")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("  Conexión perdida con el broker: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf(" Error conectando al broker: %v", token.Error())
	}

	select {}
}
