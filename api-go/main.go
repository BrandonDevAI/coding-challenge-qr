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

// RespuestaQR contiene las matrices de la factorización.
type RespuestaQR struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

func main() {
	app := crearApp()

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	log.Fatal(app.Listen(":" + puerto))
}

// crearApp arma la aplicación con sus rutas. Está separada de main
// para poder probarla en los tests sin levantar el servidor.
func crearApp() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/qr", manejarQR)

	return app
}

// manejarQR recibe la matriz, calcula su factorización QR y devuelve Q y R.
func manejarQR(c *fiber.Ctx) error {
	var solicitud SolicitudQR
	if err := c.BodyParser(&solicitud); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "el cuerpo debe ser un JSON válido"})
	}

	q, r, err := qr.Factorizar(solicitud.Matriz)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(RespuestaQR{Q: q, R: r})
}