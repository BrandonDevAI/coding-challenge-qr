package main

import (
	"log"
	"os"

	"api-go/qr"

	"github.com/gofiber/fiber/v2"
)

// SolicitudQR es el cuerpo que recibe POST /qr.
type SolicitudQR struct {
	Matriz [][]float64 `json:"matriz"`
}

// RespuestaQR contiene las matrices de la factorización
// y las estadísticas calculadas por la API Node.
type RespuestaQR struct {
	Q            [][]float64  `json:"q"`
	R            [][]float64  `json:"r"`
	Estadisticas Estadisticas `json:"estadisticas"`
}

func main() {
	puerto := leerVariable("PORT", "8080")
	urlNode := leerVariable("NODE_API_URL", "http://localhost:3005")

	app := crearApp(urlNode)
	log.Fatal(app.Listen(":" + puerto))
}

// leerVariable devuelve una variable de entorno, o el valor por defecto si no existe.
func leerVariable(nombre, porDefecto string) string {
	if valor := os.Getenv(nombre); valor != "" {
		return valor
	}
	return porDefecto
}

// crearApp arma la aplicación con sus rutas. Está separada de main
// para poder probarla en los tests sin levantar el servidor.
func crearApp(urlNode string) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/qr", manejarQR(urlNode))

	return app
}

// manejarQR recibe la matriz, calcula su factorización QR, pide las
// estadísticas a la API Node y devuelve todo junto.
func manejarQR(urlNode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var solicitud SolicitudQR
		if err := c.BodyParser(&solicitud); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "el cuerpo debe ser un JSON válido"})
		}

		q, r, err := qr.Factorizar(solicitud.Matriz)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		estadisticas, err := pedirEstadisticas(urlNode, q, r)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(RespuestaQR{Q: q, R: r, Estadisticas: estadisticas})
	}
}