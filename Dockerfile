FROM golang:1.24 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/process_inbox/main.go


FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/server .

ENTRYPOINT ["./server"]
