package domain

type DatosSensor struct {
	BME280 struct {
		Temperatura float64 `json:"temperatura"`
		Presion     float64 `json:"presion"`
		Humedad     float64 `json:"humedad"`
	} `json:"bme280"`

	MPU6050 struct {
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
	} `json:"mpu6050"`

	MLX90614 struct {
		TemperaturaAmbiente float64 `json:"temperatura_ambiente"`
		TempObjeto   float64 `json:"temp_objeto"`
	} `json:"mlx90614"`	
}
