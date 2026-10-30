const request = require('supertest');
const app = require('../src/app');

describe('GET /health', () => {
  test('responde ok', async () => {
    const res = await request(app).get('/health');

    expect(res.status).toBe(200);
    expect(res.body).toEqual({ status: 'ok' });
  });
});

describe('POST /stats', () => {
  test('devuelve las estadísticas de Q y R', async () => {
    const res = await request(app)
      .post('/stats')
      .send({ q: [[1, 0], [0, 1]], r: [[2, 3], [0, 4]] });

    expect(res.status).toBe(200);
    expect(res.body.maximo).toBe(4);
    expect(res.body.sumaTotal).toBe(11);
    expect(res.body.algunaDiagonal).toBe(true);
  });

  test('responde 400 si falta una matriz', async () => {
    const res = await request(app).post('/stats').send({ q: [[1]] });

    expect(res.status).toBe(400);
    expect(res.body.error).toBeDefined();
  });

  test('responde 400 si la matriz no es rectangular', async () => {
    const res = await request(app)
      .post('/stats')
      .send({ q: [[1, 2], [3]], r: [[1]] });

    expect(res.status).toBe(400);
  });

  test('responde 400 si no se envía cuerpo', async () => {
    const res = await request(app).post('/stats');

    expect(res.status).toBe(400);
  });
});