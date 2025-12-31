FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin ./robot_action/worker/main.go

ENTRYPOINT [ "/app/bin" ]

# TODO: Multi-stage build to reduce image size