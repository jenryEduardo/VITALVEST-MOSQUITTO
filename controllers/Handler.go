package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mosquitto/domain"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mutex        sync.Mutex
	latestSensor domain.DatosSensor

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	go func() {
		mutex.Lock()
		defer mutex.Unlock()

		fmt.Printf("📥 Mensaje recibido en el tópico '%s': %s\n", msg.Topic(), msg.Payload())

		// Procesar datos según el tópico
		switch msg.Topic() {

		// SENSOR PRINCIPAL (BME280, MPU, MLX)
		case "sensores/datos":
			var temp domain.DatosSensor
			if err := json.Unmarshal(msg.Payload(), &temp); err != nil {
				log.Printf("❌ Error al parsear sensores/datos: %v", err)
				return
			}

			// Actualiza solo lo que trae
			latestSensor.BME280 = temp.BME280
			latestSensor.MPU6050 = temp.MPU6050
			latestSensor.MLX90614 = temp.MLX90614

		// SOLO GSR
		case "GSR-SENSOR":
			var gsr struct {
				Porcentaje float64 `json:"porcentaje"`
			}
			if err := json.Unmarshal(msg.Payload(), &gsr); err != nil {
				log.Printf("❌ Error al parsear GSR: %v", err)
				return
			}

			latestSensor.GSR.Porcentaje = gsr.Porcentaje

		default:
			log.Printf("⚠️ Tópico no reconocido: %s", msg.Topic())
			return
		}

		// Convertir struct completo a JSON
		jsonData, err := json.Marshal(latestSensor)
		if err != nil {
			log.Printf("❌ Error al convertir struct a JSON: %v", err)
			return
		}

		// Endpoints donde se envían los datos
		endpoints := []string{
			"http://100.30.168.141:3000/sendData",
			"http://3.222.252.100:8081/AMQP/",
		}

		for _, url := range endpoints {
			resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("❌ Error al enviar datos a '%s': %v", url, err)
				continue
			}
			fmt.Printf("✅ Datos enviados a '%s'. Respuesta: %s\n", url, resp.Status)
			resp.Body.Close()
		}
	}()
}
