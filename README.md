# Coding Challenge: Factorización QR + Estadísticas

Dos APIs REST que trabajan juntas:

- **API Go (Fiber):** recibe una matriz, calcula su factorización QR y envía Q y R a la API Node.
- **API Node (Express):** recibe Q y R y calcula estadísticas sobre ellas.

```
Cliente ──POST /qr──▶ API Go :8080 ──POST /stats──▶ API Node :3005
        ◀──── { q, r, estadisticas } ◀───────────────────┘
```

## Cómo correrlo

### Con Docker (recomendado)

```bash
docker compose up --build
```

Esto levanta las dos APIs. Go espera a que Node esté sano (healthcheck) antes de arrancar.

### Sin Docker

Requisitos: Go 1.27+ y Node 24+.

```bash
# Terminal 1
cd api-node
npm install
npm start

# Terminal 2
cd api-go
go run .
```

## Endpoints

### `POST /qr` (API Go, puerto 8080)

Request:

```json
{ "matriz": [[12, -51, 4], [6, 167, -68], [-4, 24, -41]] }
```

Response (valores redondeados para que se lean mejor):

```json
{
  "q": [[0.857, -0.394, -0.331], [0.429, 0.903, 0.034], [-0.286, 0.171, -0.943]],
  "r": [[14, 21, -14], [0, 175, -70], [0, 0, 35]],
  "estadisticas": {
    "maximo": 175,
    "minimo": -70,
    "promedio": 8.97,
    "sumaTotal": 161.44,
    "algunaDiagonal": false
  }
}
```

Errores:

| Código | Cuándo |
|--------|--------|
| 400 | JSON inválido, matriz vacía, filas de distinto largo, más columnas que filas, o columnas linealmente dependientes |
| 502 | La API Node no responde o devuelve un error |

### `POST /stats` (API Node, puerto 3005)

Request: `{ "q": [[...]], "r": [[...]] }`. Devuelve máximo, mínimo, promedio, suma total y si alguna de las matrices es diagonal.

### `GET /health` (ambas APIs)

Devuelve `{ "status": "ok" }`. Lo usa el healthcheck de Docker.

## Tests

```bash
# API Node: Jest + Supertest (unitarios e integración)
cd api-node
npm test

# API Go: unitarios del algoritmo e integración del endpoint
cd api-go
go test ./... -v
```

Los tests de Go no necesitan la API Node: usan un servidor falso (`httptest`) que la simula, incluido el caso en que Node está caído.

## Decisiones técnicas

- **QR en vez de "rotación":** el enunciado menciona las dos cosas. Implementé la factorización QR porque es lo que pide la sección de funcionalidad requerida, y porque las estadísticas hablan de "las matrices" (Q y R).
- **Gram-Schmidt modificado:** es el método más fácil de entender y explicar. La versión modificada es más estable numéricamente que la clásica. Householder sería más robusto con matrices mal condicionadas, pero para este alcance no hacía falta. Complejidad: O(m·n²).
- **QR reducida (requiere filas ≥ columnas):** Q es de m×n y R de n×n. Si hay columnas linealmente dependientes, se responde 400 con un mensaje claro.
- **Tolerancia de 1e-10:** los cálculos con decimales dejan restos como `1e-16` donde debería haber un 0. Por eso "es diagonal" y "columna nula" se comparan con una tolerancia y no con `== 0`.
- **Lógica separada del servidor:** `qr/qr.go` y `src/estadisticas.js` no saben nada de HTTP. Así se pueden probar solas y el código queda más ordenado.
- **Estadísticas en un solo recorrido:** O(n), sin arrays intermedios. Se evita `Math.max(...array)`, que puede fallar con arrays muy grandes.
- **Comunicación Go → Node:** la URL de Node se configura con la variable `NODE_API_URL`. Así en local es `localhost:3005` y en Docker es `api-node:3005`, sin cambiar código. El cliente HTTP tiene un timeout de 5 s, y si Node falla se responde **502 Bad Gateway**, porque el error es de otro servicio y no del cliente.
- **Docker:** la imagen de Go es multi-etapa (la final pesa unos 33 MB). Ambas imágenes corren con un usuario sin privilegios, y Node instala solo las dependencias de producción.

## Estructura

```
api-go/
├── main.go            # Servidor Fiber y handler de /qr
├── cliente_node.go    # Llamada HTTP a la API Node
├── qr/qr.go           # Algoritmo de factorización QR
└── *_test.go          # Tests
api-node/
├── src/index.js         # Arranca el servidor
├── src/app.js           # Rutas de Express
├── src/estadisticas.js  # Cálculo de estadísticas
└── test/                # Tests
docker-compose.yml
```

## Pendiente / mejoras posibles

- Autenticación con JWT (opcional en el enunciado)
- Frontend para cargar matrices y ver resultados (opcional)
- Despliegue en la nube
- Factorización con Householder para matrices mal condicionadas