package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Estructura que representa los datos recibidos
type DatosSensor struct {
	Temperatura float64 `json:"temperatura"`
	Presion     float64 `json:"presion"`
	Humedad     float64 `json:"humedad"`
	Aceleracion struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"aceleracion"`
	Giroscopio struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"giroscopio"`
}

func main() {
	// Callback que se ejecuta cuando llega un mensaje
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Mensaje recibido en el topic '%s': %s\n", msg.Topic(), msg.Payload())

		var datos DatosSensor

		// Decodificar JSON en el struct
		err := json.Unmarshal(msg.Payload(), &datos)
		if err != nil {
			log.Printf("Error al parsear el JSON: %v\n", err)
			return
		}

		// Convertir el struct nuevamente a JSON para enviarlo al servidor
		jsonData, err := json.Marshal(datos)
		if err != nil {
			log.Printf("Error al convertir struct a JSON: %v\n", err)
			return
		}

		// Enviar al servidor
		resp, err := http.Post("http://localhost:8080/sendData", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error al enviar datos al servidor: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Respuesta del servidor: %s\n", resp.Status)
	}

	// Configuración del cliente MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://98.84.72.237:1883")
	opts.SetClientID("go-mqtt-client")
	opts.SetDefaultPublishHandler(messageHandler)

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
