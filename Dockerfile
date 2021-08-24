FROM node:latest AS build-frontend
WORKDIR /app
COPY frontend .
RUN npm install
RUN npm run build

FROM golang:1.16-alpine AS build-server
WORKDIR /app
COPY server .
RUN go mod download
RUN go build -o server

FROM golang:1.16-alpine
WORKDIR /app
COPY --from=build-frontend /app/dist /app/frontend/dist
COPY --from=build-server /app/server /app/server

CMD [ "/app/server",  "/app/config.yml" ]