FROM golang:1.20.10-alpine3.17 AS build-stage

LABEL authors="fadjrin"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dw-account ./cmd/main.go


FROM golang:1.20.10-alpine3.17 AS build-release-stage

WORKDIR /app

COPY --from=build-stage ./dw-account ./
COPY ./config.json ./

ENV DATABASE_MONGODB_URI=""
ENV DATABASE_MONGODB_DB_NAME=""
ENV KAFKA_BROKERS=""
ENV KAFKA_SASL_USER=""
ENV KAFKA_SASL_PASSWORD=""

RUN mkdir ./logs

EXPOSE 8000

ENTRYPOINT ["./dw-account"]