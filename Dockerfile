FROM golang:1.23 AS builder
ENV PROJECT_PATH=/build
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . ${PROJECT_PATH}
WORKDIR ${PROJECT_PATH}
RUN go build cmd/loader/main.go

FROM golang:alpine
WORKDIR /etc/synthetica
COPY --from=builder /build/main .
EXPOSE 30000
CMD ["./main"]