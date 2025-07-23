package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mosquitto/domain"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second, // Timeout para evitar bloqueos largos
}

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	// Procesamos el mensaje en una goroutine para no bloquear al cliente MQTT
	go func() {
		fmt.Printf("📩 Mensaje recibido en el topic '%s': %s\n", msg.Topic(), msg.Payload())

		var datos domain.DatosSensor

		// Decodificar el JSON en el struct
		err := json.Unmarshal(msg.Payload(), &datos)
		if err != nil {
			log.Printf("❌ Error al parsear el JSON: %v\n", err)
			return
		}

		// Convertir struct a JSON
		jsonData, err := json.Marshal(datos)
		if err != nil {
			log.Printf("❌ Error al convertir struct a JSON: %v\n", err)
			return
		}

		// Endpoints a los que se enviará el JSON
		endpoints := []string{
			"http://localhost:3000/sendData",
			"http://localhost:8085/AMQP/",
		}

		// Enviar a cada endpoint
		for _, url := range endpoints {
			resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("❌ Error al enviar datos a '%s': %v\n", url, err)
				continue
			}

			// Imprimir respuesta y cerrar el body inmediatamente
			fmt.Printf("✅ Datos enviados a '%s'. Respuesta: %s\n", url, resp.Status)
			resp.Body.Close()
		}
	}()
}
