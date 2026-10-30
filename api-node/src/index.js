const app = require('./app');

const PORT = process.env.PORT || 3005;

app.listen(PORT, () => {
  console.log(`API Node escuchando en el puerto ${PORT}`);
});