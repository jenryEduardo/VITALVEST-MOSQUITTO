package main

import (
	"mosquitto/infraestructure/controllers"
	"fmt"
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)



func main() {
	
	// Configuración del cliente MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://52.203.81.35:1883")
	opts.SetClientID("go-mqtt-client")
	opts.SetDefaultPublishHandler(controllers.MessageHandler)

	// Crear cliente
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Println("Conectado al broker Mosquitto")

	// Suscribirse al topic
	if token := client.Subscribe("sensores/datos", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Println("Suscrito al topic 'sensores/datos'")

	// Mantener el programa activo
	select {}
}
