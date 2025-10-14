ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o work-svc ./cmd/

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/config ./config
COPY --from=builder /app/work-svc ./work-svc

EXPOSE 8070

ENTRYPOINT ["./work-svc"]
