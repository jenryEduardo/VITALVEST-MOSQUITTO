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

 	// Enviar al servidor
	resp, err := http.Post("http://localhost:8080/sendData", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al enviar datos al servidor: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Respuesta del servidor: %s\n", resp.Status)
}