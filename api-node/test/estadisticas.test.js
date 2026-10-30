const {
  esMatrizValida,
  esDiagonal,
  calcularEstadisticas,
} = require('../src/estadisticas');

describe('esMatrizValida', () => {
  test('acepta una matriz rectangular de números', () => {
    expect(esMatrizValida([[1, 2], [3, 4], [5, 6]])).toBe(true);
  });

  test('rechaza filas de distinto largo', () => {
    expect(esMatrizValida([[1, 2], [3]])).toBe(false);
  });

  test('rechaza matrices vacías', () => {
    expect(esMatrizValida([])).toBe(false);
    expect(esMatrizValida([[]])).toBe(false);
  });

  test('rechaza valores que no son números', () => {
    expect(esMatrizValida([[1, 'a']])).toBe(false);
  });

  test('rechaza algo que no es un array', () => {
    expect(esMatrizValida(undefined)).toBe(false);
    expect(esMatrizValida('hola')).toBe(false);
  });
});

describe('esDiagonal', () => {
  test('la identidad es diagonal', () => {
    expect(esDiagonal([[1, 0], [0, 1]])).toBe(true);
  });

  test('una matriz con valores fuera de la diagonal no es diagonal', () => {
    expect(esDiagonal([[2, 3], [0, 4]])).toBe(false);
  });

  test('considera cero los valores muy pequeños (errores de redondeo)', () => {
    expect(esDiagonal([[1, 1e-16], [0, 1]])).toBe(true);
  });
});

describe('calcularEstadisticas', () => {
  test('calcula todas las estadísticas sobre ambas matrices', () => {
    const q = [[1, 0], [0, 1]];
    const r = [[2, 3], [0, 4]];

    expect(calcularEstadisticas([q, r])).toEqual({
      maximo: 4,
      minimo: 0,
      promedio: 1.375,
      sumaTotal: 11,
      algunaDiagonal: true,
    });
  });

  test('funciona con números negativos', () => {
    const resultado = calcularEstadisticas([[[-5, 2]], [[1, -1]]]);

    expect(resultado.maximo).toBe(2);
    expect(resultado.minimo).toBe(-5);
    expect(resultado.sumaTotal).toBe(-3);
  });
});