FROM node:24-alpine
WORKDIR /app

# Primero solo las dependencias, para aprovechar la caché de Docker
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY src ./src

# La imagen oficial de Node ya trae el usuario "node"
USER node

EXPOSE 3005
CMD ["node", "src/index.js"]