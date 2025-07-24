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
	// Variables globales para mantener los datos más recientes
	mutex        sync.Mutex
	latestSensor domain.DatosSensor
	gotSensor    bool
	gotGSR       bool

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	go func() {
		mutex.Lock()
		defer mutex.Unlock()

		fmt.Printf("📥 Mensaje recibido en el tópico '%s': %s\n", msg.Topic(), msg.Payload())

		switch msg.Topic() {
		case "sensores/datos":
			var temp domain.DatosSensor
			if err := json.Unmarshal(msg.Payload(), &temp); err != nil {
				log.Printf("❌ Error al parsear sensores/datos: %v", err)
				return
			}
			// Actualizamos solo las partes correspondientes
			latestSensor.BME280 = temp.BME280
			latestSensor.MPU6050 = temp.MPU6050
			latestSensor.MLX90614 = temp.MLX90614
			gotSensor = true

		case "GSR-SENSOR":
			var gsr struct {
				Porcentaje int `json:"porcentaje"`
			}
			if err := json.Unmarshal(msg.Payload(), &gsr); err != nil {
				log.Printf("❌ Error al parsear GSR: %v", err)
				return
			}
			latestSensor.GSR.Porcentaje = gsr.Porcentaje
			gotGSR = true

		default:
			log.Printf("⚠️ Tópico no reconocido: %s", msg.Topic())
			return
		}

		// Cuando ambos datos han sido recibidos
		if gotSensor && gotGSR {
			jsonData, err := json.Marshal(latestSensor)
			if err != nil {
				log.Printf("❌ Error al convertir struct a JSON: %v", err)
				return
			}

			endpoints := []string{
   			 "http://100.28.244.240:3000/sendData",
             "http://3.227.202.110:8085/AMQP/",
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

			// Limpiar flags para esperar nuevos datos
			gotSensor = false
			gotGSR = false
		}
	}()
}
