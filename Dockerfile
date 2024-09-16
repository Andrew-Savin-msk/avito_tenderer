FROM golang:1.22-alpine3.20

WORKDIR /zadanie

COPY . .

RUN go mod tidy

CMD go run cmd/migrator/main.go --migrations-path=migrations --direction=up && go run cmd/tenderer/main.go
