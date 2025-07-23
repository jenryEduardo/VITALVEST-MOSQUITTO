package main

import (
	"fmt"
	"log"
	"mosquitto/infraestructure/controllers"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Configuración del cliente MQTT
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://52.203.81.35:1883").
		SetClientID("go-mqtt-client").
		SetDefaultPublishHandler(controllers.MessageHandler)

	// Manejador de conexión exitosa
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Conectado al broker Mosquitto")

		// Suscribirse al topic cuando se establece la conexión
		if token := c.Subscribe("sensores/datos", 0, nil); token.Wait() && token.Error() != nil {
			log.Fatalf("❌ Error al suscribirse al topic: %v", token.Error())
		}
		fmt.Println("📡 Suscrito al topic 'sensores/datos'")
	}

	// Manejador de pérdida de conexión
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("⚠️  Conexión perdida con el broker: %v", err)
	}

	// Crear cliente
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Error conectando al broker: %v", token.Error())
	}

	// Mantener el programa activo
	select {}
}
