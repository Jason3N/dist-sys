FROM golang:1.23-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go binary (static binary for Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM scratch

COPY .env / 

COPY --from=builder /build/main /main
ENTRYPOINT ["/main"]
