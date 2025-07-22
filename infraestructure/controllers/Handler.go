package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mosquitto/domain"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Mensaje recibido en el topic '%s': %s\n", msg.Topic(), msg.Payload())

	var datos domain.DatosSensor

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

	// Lista de endpoints a los que se enviará el objeto completo
	endpoints := []string{
		"http://localhost:3000/sendData",
		"http://localhost:8085/AMQP/",
	}

	for _, url := range endpoints {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error al enviar datos a '%s': %v\n", url, err)
			continue
		}
		defer resp.Body.Close()
		fmt.Printf("✅ Datos enviados a '%s'. Respuesta: %s\n", url, resp.Status)
	}
}
