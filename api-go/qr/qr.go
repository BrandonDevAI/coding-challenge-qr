// Package qr implementa la factorización QR de una matriz usando
// el método de Gram-Schmidt modificado.
package qr

import (
	"errors"
	"math"
)

// Por debajo de este valor consideramos que una columna es cero.
const tolerancia = 1e-10

// Errores de validación. Se exportan para que quien llame pueda
// responder con un mensaje claro.
var (
	ErrMatrizVacia          = errors.New("la matriz no puede estar vacía")
	ErrNoRectangular        = errors.New("todas las filas deben tener el mismo número de columnas")
	ErrMasColumnasQueFilas  = errors.New("la matriz debe tener al menos tantas filas como columnas")
	ErrColumnasDependientes = errors.New("las columnas de la matriz son linealmente dependientes")
)

// Validar revisa que la matriz se pueda factorizar con Gram-Schmidt:
// no vacía, rectangular y con filas >= columnas.
func Validar(a [][]float64) error {
	if len(a) == 0 || len(a[0]) == 0 {
		return ErrMatrizVacia
	}

	columnas := len(a[0])
	for _, fila := range a {
		if len(fila) != columnas {
			return ErrNoRectangular
		}
	}

	if len(a) < columnas {
		return ErrMasColumnasQueFilas
	}
	return nil
}

// Factorizar descompone A (m x n) en A = Q * R, donde:
//   - Q (m x n) tiene columnas ortonormales (perpendiculares y de largo 1)
//   - R (n x n) es triangular superior
func Factorizar(a [][]float64) ([][]float64, [][]float64, error) {
	if err := Validar(a); err != nil {
		return nil, nil, err
	}

	m, n := len(a), len(a[0])

	// Copiamos A en Q para no modificar la matriz original.
	q := nuevaMatriz(m, n)
	for i := range a {
		copy(q[i], a[i])
	}
	r := nuevaMatriz(n, n)

	for k := 0; k < n; k++ {
		// 1. Calculamos el largo (norma) de la columna k.
		norma := 0.0
		for i := 0; i < m; i++ {
			norma += q[i][k] * q[i][k]
		}
		norma = math.Sqrt(norma)

		if norma < tolerancia {
			return nil, nil, ErrColumnasDependientes
		}

		// 2. Normalizamos la columna k para que mida 1.
		r[k][k] = norma
		for i := 0; i < m; i++ {
			q[i][k] /= norma
		}

		// 3. A cada columna siguiente le quitamos la parte que
		//    apunta en la misma dirección que la columna k.
		for j := k + 1; j < n; j++ {
			producto := 0.0
			for i := 0; i < m; i++ {
				producto += q[i][k] * q[i][j]
			}

			r[k][j] = producto
			for i := 0; i < m; i++ {
				q[i][j] -= producto * q[i][k]
			}
		}
	}

	return q, r, nil
}

// nuevaMatriz crea una matriz de ceros con el tamaño indicado.
func nuevaMatriz(filas, columnas int) [][]float64 {
	matriz := make([][]float64, filas)
	for i := range matriz {
		matriz[i] = make([]float64, columnas)
	}
	return matriz
}