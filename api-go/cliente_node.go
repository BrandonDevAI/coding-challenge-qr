package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Estadisticas es lo que devuelve la API Node en POST /stats.
type Estadisticas struct {
	Maximo         float64 `json:"maximo"`
	Minimo         float64 `json:"minimo"`
	Promedio       float64 `json:"promedio"`
	SumaTotal      float64 `json:"sumaTotal"`
	AlgunaDiagonal bool    `json:"algunaDiagonal"`
}

// Cliente HTTP con timeout, para no quedarnos esperando para siempre
// si la API Node no responde. Es generoso porque en el plan gratuito
// de Render el servicio "despierta" en hasta ~1 minuto.
var clienteHTTP = &http.Client{Timeout: 60 * time.Second}

func pedirEstadisticas(urlNode string, q, r [][]float64) (Estadisticas, error) {
	var estadisticas Estadisticas

	cuerpo, err := json.Marshal(map[string][][]float64{"q": q, "r": r})
	if err != nil {
		return estadisticas, err
	}

	resp, err := clienteHTTP.Post(urlNode+"/stats", "application/json", bytes.NewReader(cuerpo))
	if err != nil {
		return estadisticas, fmt.Errorf("no se pudo conectar con la API Node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return estadisticas, fmt.Errorf("la API Node respondió con estado %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&estadisticas)
	return estadisticas, err
}