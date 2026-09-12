FROM docker.io/library/golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o devops-homework ./cmd

FROM gcr.io/distroless/base:nonroot

WORKDIR /app

COPY --from=builder /app/devops-homework .

EXPOSE 8080

CMD ["./devops-homework"]