package domain

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
