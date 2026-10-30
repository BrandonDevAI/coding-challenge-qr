const express = require('express');
const { esMatrizValida, calcularEstadisticas } = require('./estadisticas');

const app = express();
app.use(express.json());

// Endpoint de salud: sirve para probar rápido que la API está arriba (y luego para Docker)
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

/**
 * POST /stats
 * Body esperado: { "q": [[...]], "r": [[...]] }
 * Recibe las matrices Q y R que envía la API de Go y devuelve sus estadísticas.
 */
app.post('/stats', (req, res) => {
  const { q, r } = req.body ?? {};

  if (!esMatrizValida(q) || !esMatrizValida(r)) {
    return res
      .status(400)
      .json({ error: 'q y r deben ser matrices rectangulares de números' });
  }

  res.json(calcularEstadisticas([q, r]));
});

module.exports = app;