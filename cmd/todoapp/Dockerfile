FROM golang:1.27.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.23
RUN addgroup -S todoapp && adduser -S -G todoapp todoapp \
	&& mkdir -p /app/out/logs && chown -R todoapp:todoapp /app
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app
USER todoapp
ENV HTTP_ADDR=:5050
EXPOSE 5050
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
	CMD wget -qO- http://127.0.0.1:5050/health/live || exit 1
CMD ["/app/exe"]
