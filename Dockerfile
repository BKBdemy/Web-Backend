FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .

# set up mockery & swaggo
RUN go install github.com/vektra/mockery/v2@latest
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

RUN go mod download
RUN go generate ./...
RUN swag init

RUN go build -o bkbdemy-backend EntitlementServer

FROM alpine:latest
RUN apk --no-cache add ca-certificates gcompat && update-ca-certificates

COPY --from=builder /app/bkbdemy-backend /app/bkbdemy-backend
EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/bkbdemy-backend"]