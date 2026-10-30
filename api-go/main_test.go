package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// nodeFalso levanta un servidor de prueba que se hace pasar por la API Node.
func nodeFalso(t *testing.T) *httptest.Server {
	t.Helper()

	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stats" {
			t.Errorf("Go llamó a %s, esperaba /stats", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"maximo":2,"minimo":0,"promedio":0.75,"sumaTotal":6,"algunaDiagonal":true}`))
	}))
	t.Cleanup(servidor.Close)

	return servidor
}

// enviarQR hace un POST /qr a la app con el cuerpo indicado.
func enviarQR(t *testing.T, urlNode, cuerpo string) *http.Response {
	t.Helper()

	app := crearApp(urlNode)
	req := httptest.NewRequest(http.MethodPost, "/qr", strings.NewReader(cuerpo))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("error al hacer el request: %v", err)
	}
	return resp
}

func TestQRDevuelveMatricesYEstadisticas(t *testing.T) {
	node := nodeFalso(t)

	resp := enviarQR(t, node.URL, `{"matriz":[[1,0],[0,2]]}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("estado = %d, esperaba 200", resp.StatusCode)
	}

	var respuesta RespuestaQR
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		t.Fatalf("no se pudo leer la respuesta: %v", err)
	}

	if len(respuesta.Q) != 2 || len(respuesta.R) != 2 {
		t.Errorf("esperaba Q y R de 2x2, obtuve %v y %v", respuesta.Q, respuesta.R)
	}
	if !respuesta.Estadisticas.AlgunaDiagonal || respuesta.Estadisticas.SumaTotal != 6 {
		t.Errorf("las estadísticas no son las que mandó Node: %+v", respuesta.Estadisticas)
	}
}

func TestQRConMatrizInvalidaDevuelve400(t *testing.T) {
	node := nodeFalso(t)

	resp := enviarQR(t, node.URL, `{"matriz":[[1,2],[3]]}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("estado = %d, esperaba 400", resp.StatusCode)
	}
}

func TestQRConJSONInvalidoDevuelve400(t *testing.T) {
	node := nodeFalso(t)

	resp := enviarQR(t, node.URL, `esto no es json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("estado = %d, esperaba 400", resp.StatusCode)
	}
}

func TestQRSinNodeDevuelve502(t *testing.T) {
	// Levantamos un servidor y lo cerramos enseguida: así tenemos
	// una URL donde seguro no hay nadie escuchando.
	node := httptest.NewServer(http.NotFoundHandler())
	node.Close()

	resp := enviarQR(t, node.URL, `{"matriz":[[1,0],[0,2]]}`)
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("estado = %d, esperaba 502", resp.StatusCode)
	}
}