// Tolerancia para comparar con cero.
// La factorización QR trabaja con decimales, así que un valor como 1e-16
// en realidad es un cero con error de redondeo.
const TOLERANCIA = 1e-10;

/**
 * Valida que sea una matriz rectangular de números:
 * un array no vacío de filas, todas del mismo largo.
 * @param {any} matriz
 * @returns {boolean}
 */
function esMatrizValida(matriz) {
  if (!Array.isArray(matriz) || matriz.length === 0) return false;
  if (!Array.isArray(matriz[0]) || matriz[0].length === 0) return false;

  const columnas = matriz[0].length;

  return matriz.every(
    (fila) =>
      Array.isArray(fila) &&
      fila.length === columnas &&
      fila.every((valor) => typeof valor === 'number' && Number.isFinite(valor))
  );
}

/**
 * Indica si una matriz es diagonal: todos los valores fuera
 * de la diagonal principal son (casi) cero.
 * @param {number[][]} matriz
 * @returns {boolean}
 */
function esDiagonal(matriz) {
  for (let i = 0; i < matriz.length; i++) {
    for (let j = 0; j < matriz[i].length; j++) {
      if (i !== j && Math.abs(matriz[i][j]) > TOLERANCIA) {
        return false;
      }
    }
  }
  return true;
}

/**
 * Calcula máximo, mínimo, promedio, suma total y si alguna
 * matriz es diagonal. Recorre todos los valores una sola vez.
 * @param {number[][][]} matrices
 */
function calcularEstadisticas(matrices) {
  let maximo = -Infinity;
  let minimo = Infinity;
  let suma = 0;
  let cantidad = 0;

  for (const matriz of matrices) {
    for (const fila of matriz) {
      for (const valor of fila) {
        if (valor > maximo) maximo = valor;
        if (valor < minimo) minimo = valor;
        suma += valor;
        cantidad++;
      }
    }
  }

  return {
    maximo,
    minimo,
    promedio: suma / cantidad,
    sumaTotal: suma,
    algunaDiagonal: matrices.some(esDiagonal),
  };
}

module.exports = { esMatrizValida, esDiagonal, calcularEstadisticas };