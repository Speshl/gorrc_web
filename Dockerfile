# syntax=docker/dockerfile:1

#Command: docker build -t gorrc_web -f Dockerfile .

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY internal/ ./internal
COPY public/ ./public

RUN CGO_ENABLED=0 GOOS=linux go build -o /gorrc_web

CMD ["/gorrc_web"]