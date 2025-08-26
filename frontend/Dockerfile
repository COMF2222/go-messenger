FROM node:20 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:latest AS frontend
COPY --from=builder /app/build /usr/share/nginx/html
