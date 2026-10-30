package qr

import (
	"errors"
	"math"
	"testing"
)

// casiIgual compara dos números aceptando un pequeño error de redondeo.
func casiIgual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestFactorizarEjemploClasico(t *testing.T) {
	a := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}
	rEsperada := [][]float64{
		{14, 21, -14},
		{0, 175, -70},
		{0, 0, 35},
	}

	_, r, err := Factorizar(a)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}

	for i := range rEsperada {
		for j := range rEsperada[i] {
			if !casiIgual(r[i][j], rEsperada[i][j]) {
				t.Errorf("R[%d][%d] = %v, esperaba %v", i, j, r[i][j], rEsperada[i][j])
			}
		}
	}
}

func TestQPorRDevuelveLaMatrizOriginal(t *testing.T) {
	// Matriz rectangular (4 x 2) para probar también el caso no cuadrado.
	a := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 9},
	}

	q, r, err := Factorizar(a)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}

	// Multiplicamos Q * R y comparamos con A.
	for i := range a {
		for j := range a[i] {
			suma := 0.0
			for k := range r {
				suma += q[i][k] * r[k][j]
			}
			if !casiIgual(suma, a[i][j]) {
				t.Errorf("(Q*R)[%d][%d] = %v, esperaba %v", i, j, suma, a[i][j])
			}
		}
	}
}

func TestColumnasDeQSonOrtonormales(t *testing.T) {
	a := [][]float64{
		{2, -1, 0},
		{1, 3, 1},
		{0, 1, 4},
	}

	q, _, err := Factorizar(a)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}

	// El producto punto de dos columnas debe ser 1 si es la misma columna
	// (largo 1) y 0 si son distintas (perpendiculares).
	for c1 := 0; c1 < 3; c1++ {
		for c2 := 0; c2 < 3; c2++ {
			producto := 0.0
			for i := range q {
				producto += q[i][c1] * q[i][c2]
			}

			esperado := 0.0
			if c1 == c2 {
				esperado = 1
			}
			if !casiIgual(producto, esperado) {
				t.Errorf("columna %d · columna %d = %v, esperaba %v", c1, c2, producto, esperado)
			}
		}
	}
}

func TestFactorizarErrores(t *testing.T) {
	casos := []struct {
		nombre   string
		matriz   [][]float64
		esperado error
	}{
		{"matriz vacía", [][]float64{}, ErrMatrizVacia},
		{"no rectangular", [][]float64{{1, 2}, {3}}, ErrNoRectangular},
		{"más columnas que filas", [][]float64{{1, 2, 3}}, ErrMasColumnasQueFilas},
		{"columnas dependientes", [][]float64{{1, 2}, {2, 4}}, ErrColumnasDependientes},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			_, _, err := Factorizar(caso.matriz)
			if !errors.Is(err, caso.esperado) {
				t.Errorf("esperaba error %q, obtuve %v", caso.esperado, err)
			}
		})
	}
}